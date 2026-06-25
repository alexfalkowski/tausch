package io

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
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
// decoded before being copied to w. Relative file payloads are resolved from dir.
// On successful output, Write returns (true, nil).
//
// The returned boolean indicates whether output was attempted/emitted. This is used
// by the CLI orchestration to decide whether stdout was produced or whether it
// should fall back to stderr.
//
// Errors from decoding or from the underlying writer are returned.
func Write(w io.Writer, data, dir string) (bool, error) {
	if data == "" {
		return false, nil
	}

	r, err := decode(data, dir)
	if err != nil {
		return false, err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return true, err
}

func decode(value, dir string) (io.ReadCloser, error) {
	kind, data, ok := strings.Cut(value, ":")
	if !ok {
		return nil, ErrKindNotFound
	}

	switch kind {
	case "text":
		return io.NopCloser(strings.NewReader(data)), nil
	case "file":
		return os.Open(filePath(data, dir))
	case "base64":
		d, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}

		return io.NopCloser(bytes.NewReader(d)), nil
	}

	return nil, ErrKindNotFound
}

func filePath(path, dir string) string {
	if dir == "" || filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(dir, path)
}
