package encoding_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/tausch/internal/encoding"
	"github.com/stretchr/testify/require"
)

func TestDecodeSuccess(t *testing.T) {
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
			d, err := encoding.Decode(tt.value)

			require.NoError(t, err)
			require.Equal(t, want, d)
		})
	}
}

func TestDecodeError(t *testing.T) {
	values := []string{
		"bob:test",
		"base64:1234567",
		"file:../../test/configs/none.txt",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			_, err := encoding.Decode(value)
			require.Error(t, err)
		})
	}
}

func TestDecodeMissingSeparator(t *testing.T) {
	values := []string{
		"text",
		"base64",
		"file",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			_, err := encoding.Decode(value)
			require.ErrorIs(t, err, encoding.ErrKindNotFound)
		})
	}
}

func TestDecodePreservesDataAfterFirstColon(t *testing.T) {
	d, err := encoding.Decode("text:a:b")

	require.NoError(t, err)
	require.Equal(t, []byte("a:b"), d)
}
