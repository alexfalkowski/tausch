package io_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/io"
	"github.com/stretchr/testify/assert"
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

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, []byte("test"), bytes.TrimSpace(buffer.Bytes()))
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

		assert.Error(t, err)
		assert.False(t, ok)
		assert.Empty(t, buffer.Bytes())
	}
}

func TestWriteEmpty(t *testing.T) {
	buffer := &bytes.Buffer{}
	ok, err := io.Write(buffer, "")

	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Empty(t, buffer.Bytes())
}
