package encoding

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

// ErrKindNotFound is returned when a value cannot be decoded because its kind
// prefix is not supported.
//
// A "kind" is the substring before the first ":" in a tausch payload string.
// For example, in "text:hello", the kind is "text".
var ErrKindNotFound = errors.New("kind not found")

// Decode converts a tausch payload string in the form "kind:data" into raw bytes.
//
// The kind is determined by taking the substring before the first ":" character.
// Everything after the first ":" is treated as the data portion (and may itself
// contain additional ":" characters).
//
// Supported kinds:
//
//   - "text": returns the bytes of data as-is.
//   - "file": treats data as a filesystem path and returns the file contents.
//   - "base64": treats data as standard base64 and decodes it.
//
// If kind is unknown, Decode returns ErrKindNotFound.
//
// If kind is "file", any error from os.ReadFile is returned.
// If kind is "base64", any error from base64.StdEncoding.DecodeString is returned.
func Decode(value string) ([]byte, error) {
	kind, data, _ := strings.Cut(value, ":")
	switch kind {
	case "text":
		return []byte(data), nil
	case "file":
		return os.ReadFile(data)
	case "base64":
		return base64.StdEncoding.DecodeString(data)
	}

	return nil, ErrKindNotFound
}
