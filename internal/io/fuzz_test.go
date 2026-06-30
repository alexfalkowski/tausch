package io_test

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/alexfalkowski/tausch/internal/io"
	"github.com/stretchr/testify/require"
)

// FuzzWriteText protects user-controlled text payloads, where configured
// command output can be empty or contain separators after the kind prefix.
func FuzzWriteText(f *testing.F) {
	f.Add("")
	f.Add(" test\n")
	f.Add("a:b")

	f.Fuzz(func(t *testing.T, data string) {
		if len(data) > 4096 {
			t.Skip()
		}

		buffer := &bytes.Buffer{}
		wrote, err := io.Write(buffer, "text:"+data, "")

		require.NoError(t, err)
		require.True(t, wrote)
		require.Equal(t, data, buffer.String())
	})
}

// FuzzWriteBase64 protects base64 payload decoding for arbitrary recorded
// command output bytes.
func FuzzWriteBase64(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte(" test\n"))
	f.Add([]byte{0x00, 0x01, 0xff})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			t.Skip()
		}

		buffer := &bytes.Buffer{}
		wrote, err := io.Write(buffer, "base64:"+base64.StdEncoding.EncodeToString(data), "")

		require.NoError(t, err)
		require.True(t, wrote)
		require.Equal(t, string(data), buffer.String())
	})
}

// FuzzWriteKindNotFound protects the kind:data parser from accepting unknown
// output encodings.
func FuzzWriteKindNotFound(f *testing.F) {
	f.Add("bob", "test")
	f.Add("", "test")
	f.Add("unsupported", "")

	f.Fuzz(func(t *testing.T, kind, data string) {
		if len(kind)+len(data) > 4096 {
			t.Skip()
		}

		value := kind + ":" + data
		prefix, _, _ := strings.Cut(value, ":")
		if prefix == "text" || prefix == "file" || prefix == "base64" {
			t.Skip()
		}

		buffer := &bytes.Buffer{}
		wrote, err := io.Write(buffer, value, "")

		require.ErrorIs(t, err, io.ErrKindNotFound)
		require.False(t, wrote)
		require.Empty(t, buffer.Bytes())
	})
}
