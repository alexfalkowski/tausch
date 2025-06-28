package io

import (
	"io"

	"github.com/alexfalkowski/tausch/internal/encoding"
)

// Writer is an alias for io.Writer.
type Writer = io.Writer

// Write data to w.
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
