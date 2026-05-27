package test

import "errors"

// ErrWriteFailed is returned by FailingWriter.
var ErrWriteFailed = errors.New("write failed")

// FailingWriter is an io.Writer that always fails.
type FailingWriter struct{}

func (FailingWriter) Write([]byte) (int, error) {
	return 0, ErrWriteFailed
}
