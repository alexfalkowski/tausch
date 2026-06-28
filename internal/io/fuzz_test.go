package io_test

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/alexfalkowski/tausch/internal/io"
	"github.com/stretchr/testify/require"
)

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
