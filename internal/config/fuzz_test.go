package config_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/stretchr/testify/require"
)

// FuzzConfigValidate protects YAML-decoded command invariants from arbitrary
// command names, output payloads, and exit code values.
func FuzzConfigValidate(f *testing.F) {
	f.Add("go version", "text:go version", "", 0, false)
	f.Add("go bob", "", "text:go bob", 127, true)
	f.Add("go version", "text:go version", "text:warning", 0, false)
	f.Add("go bob", "", "", -1, true)
	f.Add("go bob", "", "", 256, true)

	f.Fuzz(func(t *testing.T, name, stdout, stderr string, exitCode int, hasExitCode bool) {
		if len(name)+len(stdout)+len(stderr) > 4096 {
			t.Skip()
		}

		var code *int
		if hasExitCode {
			code = &exitCode
		}

		c := &config.Config{
			Cmds: []*config.Command{
				{Name: name, Stdout: stdout, Stderr: stderr, ExitCode: code},
			},
		}

		err := c.Validate()
		switch {
		case stdout != "" && stderr != "":
			require.ErrorIs(t, err, config.ErrMultipleOutputs)
		case hasExitCode && (exitCode < 0 || exitCode > 255):
			require.ErrorIs(t, err, config.ErrInvalidExitCode)
		default:
			require.NoError(t, err)
		}
	})
}

// FuzzConfigGetCommand protects exact command-name matching, which is the CLI's
// user-visible lookup contract.
func FuzzConfigGetCommand(f *testing.F) {
	f.Add("go version", "go version")
	f.Add("go version", "Go Version")
	f.Add("go version", " go version ")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, name, query string) {
		if len(name)+len(query) > 4096 {
			t.Skip()
		}

		command := &config.Command{Name: name}
		c := &config.Config{Cmds: []*config.Command{nil, command}}

		got, err := c.GetCommand(query)
		if query == name {
			require.NoError(t, err)
			require.Same(t, command, got)

			return
		}

		require.Nil(t, got)
		require.ErrorIs(t, err, config.ErrCommandNotFound)
	})
}
