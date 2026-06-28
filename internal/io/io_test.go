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
	t.Parallel()

	want := []byte(" test\n")
	file := filepath.Join(t.TempDir(), "sample.txt")
	require.NoError(t, os.WriteFile(file, want, 0o600))

	values := []struct {
		dir   string
		name  string
		value string
	}{
		{name: "text", value: "text: test\n"},
		{name: "base64", value: "base64:IHRlc3QK"},
		{name: "file absolute", value: "file:" + file},
		{dir: filepath.Dir(file), name: "file relative", value: "file:" + filepath.Base(file)},
	}

	for _, tt := range values {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buffer := &bytes.Buffer{}
			wrote, err := io.Write(buffer, tt.value, tt.dir)

			require.NoError(t, err)
			require.True(t, wrote)
			require.Equal(t, want, buffer.Bytes())
		})
	}
}

func TestWriteError(t *testing.T) {
	t.Parallel()

	values := []string{
		"base64:123456",
		"file:../../test/configs/none.txt",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			buffer := &bytes.Buffer{}
			wrote, err := io.Write(buffer, value, "")

			require.Error(t, err)
			require.False(t, wrote)
			require.Empty(t, buffer.Bytes())
		})
	}
}

func TestWriteKindNotFound(t *testing.T) {
	t.Parallel()

	values := []string{
		"bob:test",
		"text",
		"base64",
		"file",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			buffer := &bytes.Buffer{}
			wrote, err := io.Write(buffer, value, "")

			require.ErrorIs(t, err, io.ErrKindNotFound)
			require.False(t, wrote)
			require.Empty(t, buffer.Bytes())
		})
	}
}

func TestWritePreservesDataAfterFirstColon(t *testing.T) {
	t.Parallel()

	buffer := &bytes.Buffer{}
	wrote, err := io.Write(buffer, "text:a:b", "")

	require.NoError(t, err)
	require.True(t, wrote)
	require.Equal(t, []byte("a:b"), buffer.Bytes())
}

func TestWriteEmpty(t *testing.T) {
	t.Parallel()

	buffer := &bytes.Buffer{}
	wrote, err := io.Write(buffer, "", "")

	require.NoError(t, err)
	require.False(t, wrote)
	require.Empty(t, buffer.Bytes())
}

func TestWriteWriterError(t *testing.T) {
	t.Parallel()

	wrote, err := io.Write(test.FailingWriter{}, "text:test", "")

	require.ErrorIs(t, err, test.ErrWriteFailed)
	require.True(t, wrote)
}
