package main

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/cmd"
)

// BenchmarkRunStdout tracks the CLI's hot stub path: parse argv, decode YAML,
// find the command, decode stdout, and write the configured bytes.
func BenchmarkRunStdout(b *testing.B) {
	args := []string{
		"-config", "test/configs/config.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		stdout.Reset()
		stderr.Reset()

		code := cmd.Run(stdout, stderr, args)
		if code != 0 {
			b.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
		}
		if stdout.Len() == 0 {
			b.Fatal("expected stdout")
		}
	}
}
