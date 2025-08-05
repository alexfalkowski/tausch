package encoding

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

// ErrKindNotFound for encoding.
var ErrKindNotFound = errors.New("kind not found")

// Decode the value, which is kind:data.
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
