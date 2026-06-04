package io_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/tausch/internal/io"
	"github.com/alexfalkowski/tausch/internal/test"
	"github.com/stretchr/testify/require"
)

func TestWriteSuccess(t *testing.T) {
	want := []byte(" test\n")
	file := filepath.Join(t.TempDir(), "sample.txt")
	require.NoError(t, os.WriteFile(file, want, 0o600))

	values := []struct {
		name  string
		value string
	}{
		{name: "text", value: "text: test\n"},
		{name: "base64", value: "base64:IHRlc3QK"},
		{name: "file", value: "file:" + file},
	}

	for _, tt := range values {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			wrote, err := io.Write(buffer, tt.value)

			require.NoError(t, err)
			require.True(t, wrote)
			require.Equal(t, want, buffer.Bytes())
		})
	}
}

func TestWriteError(t *testing.T) {
	values := []string{
		"base64:123456",
		"file:../../test/configs/none.txt",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			wrote, err := io.Write(buffer, value)

			require.Error(t, err)
			require.False(t, wrote)
			require.Empty(t, buffer.Bytes())
		})
	}
}

func TestWriteEmpty(t *testing.T) {
	buffer := &bytes.Buffer{}
	wrote, err := io.Write(buffer, "")

	require.NoError(t, err)
	require.False(t, wrote)
	require.Empty(t, buffer.Bytes())
}

func TestWriteWriterError(t *testing.T) {
	wrote, err := io.Write(test.FailingWriter{}, "text:test")

	require.ErrorIs(t, err, test.ErrWriteFailed)
	require.True(t, wrote)
}
