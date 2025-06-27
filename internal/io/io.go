package io

import (
	"io"

	"github.com/alexfalkowski/tausch/internal/encoding"
)

// Write data to w.
func Write(data string, w io.Writer) (bool, error) {
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
