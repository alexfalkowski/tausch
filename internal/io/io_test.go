package io_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/io"
	"github.com/stretchr/testify/require"
)

func TestWriteSuccess(t *testing.T) {
	values := []string{
		"text:test",
		"base64:dGVzdA==",
		"file:../../test/configs/test.txt",
	}

	for _, value := range values {
		buffer := &bytes.Buffer{}
		ok, err := io.Write(buffer, value)

		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, []byte("test"), bytes.TrimSpace(buffer.Bytes()))
	}
}

func TestWriteError(t *testing.T) {
	values := []string{
		"base64:123456",
		"file:../../test/configs/none.txt",
	}

	for _, value := range values {
		buffer := &bytes.Buffer{}
		ok, err := io.Write(buffer, value)

		require.Error(t, err)
		require.False(t, ok)
		require.Empty(t, buffer.Bytes())
	}
}

func TestWriteEmpty(t *testing.T) {
	buffer := &bytes.Buffer{}
	ok, err := io.Write(buffer, "")

	require.NoError(t, err)
	require.False(t, ok)
	require.Empty(t, buffer.Bytes())
}
