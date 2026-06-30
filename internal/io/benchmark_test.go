package io_test

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/alexfalkowski/tausch/internal/io"
)

// BenchmarkWriteText tracks text payload decoding and write throughput for
// recorded command output.
func BenchmarkWriteText(b *testing.B) {
	value := "text:" + strings.Repeat("x", 1024)
	buffer := &bytes.Buffer{}

	b.ReportAllocs()
	b.SetBytes(1024)
	b.ResetTimer()

	for range b.N {
		buffer.Reset()

		wrote, err := io.Write(buffer, value, "")
		if err != nil {
			b.Fatal(err)
		}
		if !wrote {
			b.Fatal("expected write")
		}
	}
}

// BenchmarkWriteBase64 tracks base64 payload decoding and write throughput for
// recorded command output bytes.
func BenchmarkWriteBase64(b *testing.B) {
	data := []byte(strings.Repeat("x", 1024))
	value := "base64:" + base64.StdEncoding.EncodeToString(data)
	buffer := &bytes.Buffer{}

	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for range b.N {
		buffer.Reset()

		wrote, err := io.Write(buffer, value, "")
		if err != nil {
			b.Fatal(err)
		}
		if !wrote {
			b.Fatal("expected write")
		}
	}
}
