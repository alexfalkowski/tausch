package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/alexfalkowski/tausch/internal/cmd"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func FuzzRunWritesConfiguredStdout(f *testing.F) {
	f.Add("go version")
	f.Add(" go version ")
	f.Add("-version --json")

	f.Fuzz(func(t *testing.T, rawName string) {
		if len(rawName) > 1024 || !utf8.ValidString(rawName) || strings.ContainsRune(rawName, 0) {
			t.Skip()
		}

		name := strings.TrimSpace(rawName)
		if name == "" {
			t.Skip()
		}

		data, err := yaml.Marshal(fuzzConfig{
			Cmds: []fuzzCommand{
				{Name: name, Stdout: "text:ok"},
			},
		})
		require.NoError(t, err)

		config := filepath.Join(t.TempDir(), "config.yml")
		require.NoError(t, os.WriteFile(config, data, 0o600))

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		code := cmd.Run(stdout, stderr, []string{"-config", config, "--", rawName})

		require.Zero(t, code)
		require.Equal(t, "ok", stdout.String())
		require.Empty(t, stderr.String())
	})
}

type fuzzConfig struct {
	Cmds []fuzzCommand `yaml:"cmds"`
}

type fuzzCommand struct {
	Name   string `yaml:"name"`
	Stdout string `yaml:"stdout"`
}
