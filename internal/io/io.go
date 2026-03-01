package io

import (
	"io"

	"github.com/alexfalkowski/tausch/internal/encoding"
)

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
// If data is non-empty, it is expected to be in tausch's `kind:data` format and is
// decoded via [encoding.Decode] before being written to w. On successful decode and
// write, Write returns (true, nil).
//
// The returned boolean indicates whether output was attempted/emitted. This is used
// by the CLI orchestration to decide whether a command should be treated as having
// produced stdout (success path) or should fall back to stderr (error path).
//
// Errors from decoding or from the underlying writer are returned.
func Write(w io.Writer, data string) (bool, error) {
	if len(data) == 0 {
		return false, nil
	}

	d, err := encoding.Decode(data)
	if err != nil {
		return false, err
	}

	_, err = w.Write(d)
	return true, err
}
