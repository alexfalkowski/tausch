package cmd_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/cmd"
	"github.com/stretchr/testify/require"
)

func TestRunInvalidArgs(t *testing.T) {
	args := []string{"- x"}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "flag provided but not defined")
	require.Contains(t, stderr.String(), "Usage of tausch:")
}

func TestRunConfigError(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, nil)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.NotEmpty(t, stderr.Bytes())
}

func TestRunMissingConfig(t *testing.T) {
	args := []string{
		"-config", "cfg.yml",
		"--",
		"test", "my", "code",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "cfg.yml")
}

func TestRunMissingCommand(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"test", "my", "code",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "find test my code: command not found")
}

func TestRunStdoutWriteError(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/stdout_invalid.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "illegal base64 data")
}

func TestRunStdout(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Zero(t, c)
	require.NotEmpty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunStderrWriteError(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/stderr_invalid.yml",
		"--",
		"go", "bob",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "illegal base64 data")
}

func TestRunStderr(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"go", "bob",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.NotEmpty(t, stderr.Bytes())
}
