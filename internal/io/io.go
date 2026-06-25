package io

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
)

// ErrKindNotFound is returned when a value cannot be decoded because its kind
// prefix is not supported.
var ErrKindNotFound = errors.New("kind not found")

// Writer is an alias for [io.Writer].
//
// It is provided as a small convenience so callers working primarily with
// tausch's internal packages don't also need to import the standard library io
// package just to name the interface type.
type Writer = io.Writer

// Write decodes data and writes it to w.
//
// If data is the empty string, Write performs no writes and returns (false, nil).
//
// If data is non-empty, it is expected to be in tausch's `kind:data` format and
// decoded before being copied to w. On successful output, Write returns (true, nil).
//
// The returned boolean indicates whether output was attempted/emitted. This is used
// by the CLI orchestration to decide whether a command should be treated as having
// produced stdout (success path) or should fall back to stderr (error path).
//
// Errors from decoding or from the underlying writer are returned.
func Write(w io.Writer, data string) (bool, error) {
	if data == "" {
		return false, nil
	}

	r, err := decode(data)
	if err != nil {
		return false, err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return true, err
}

func decode(value string) (io.ReadCloser, error) {
	kind, data, ok := strings.Cut(value, ":")
	if !ok {
		return nil, ErrKindNotFound
	}

	switch kind {
	case "text":
		return io.NopCloser(strings.NewReader(data)), nil
	case "file":
		return os.Open(data)
	case "base64":
		d, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}

		return io.NopCloser(bytes.NewReader(d)), nil
	}

	return nil, ErrKindNotFound
}
