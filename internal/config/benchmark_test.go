package config_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/config"
)

// BenchmarkDecode tracks the cost of loading and validating the YAML config on
// each CLI invocation.
func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		c, err := config.Decode("../../test/configs/config.yml")
		if err != nil {
			b.Fatal(err)
		}
		if c == nil {
			b.Fatal("expected config")
		}
	}
}

// BenchmarkGetCommand tracks exact command lookup cost as config entries grow.
func BenchmarkGetCommand(b *testing.B) {
	command := &config.Command{Name: "go version"}
	c := &config.Config{
		Cmds: []*config.Command{
			nil,
			{Name: "go env GOPATH"},
			{Name: "go test ./..."},
			command,
			{Name: "go build"},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		got, err := c.GetCommand("go version")
		if err != nil {
			b.Fatal(err)
		}
		if got != command {
			b.Fatal("expected matching command")
		}
	}
}
