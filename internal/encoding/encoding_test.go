package encoding_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/encoding"
	"github.com/stretchr/testify/assert"
)

func TestDecodeSuccess(t *testing.T) {
	values := []string{
		"text:test",
		"base64:dGVzdA==",
		"file:../../test/configs/test.txt",
	}

	for _, value := range values {
		d, err := encoding.Decode(value)

		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), bytes.TrimSpace(d))
	}
}

func TestDecodeError(t *testing.T) {
	values := []string{
		"bob:test",
		"base64:1234567",
		"file:../../test/configs/none.txt",
	}

	for _, value := range values {
		_, err := encoding.Decode(value)
		assert.Error(t, err)
	}
}
