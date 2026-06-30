package flag_test

import (
	"io"
	"testing"

	"github.com/alexfalkowski/tausch/internal/flag"
)

// BenchmarkParseCommandName tracks the CLI argument parsing and command-name
// derivation performed before every config lookup.
func BenchmarkParseCommandName(b *testing.B) {
	args := []string{"-config", "config.yml", "--", "go", "version"}

	b.ReportAllocs()

	for range b.N {
		values, err := flag.Parse(io.Discard, args)
		if err != nil {
			b.Fatal(err)
		}
		if values.Name() != "go version" {
			b.Fatalf("expected go version, got %q", values.Name())
		}
	}
}
