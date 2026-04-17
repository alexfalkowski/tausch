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
	c, err := cmd.Run(stdout, stderr, args)

	require.Error(t, err)
	require.Zero(t, c)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "flag provided but not defined")
	require.Contains(t, stderr.String(), "Usage of tausch:")
}

func TestRunMissingConfig(t *testing.T) {
	args := []string{
		"-config", "cfg.yml",
		"--",
		"test", "my", "code",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c, err := cmd.Run(stdout, stderr, args)

	require.Error(t, err)
	require.Zero(t, c)
	require.Empty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunMissingCommand(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"test", "my", "code",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c, err := cmd.Run(stdout, stderr, args)

	require.Error(t, err)
	require.Zero(t, c)
	require.Empty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunStdout(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c, err := cmd.Run(stdout, stderr, args)

	require.NoError(t, err)
	require.Zero(t, c)
	require.NotEmpty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunStderr(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"go", "bob",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c, err := cmd.Run(stdout, stderr, args)

	require.NoError(t, err)
	require.Equal(t, 1, c)
	require.Empty(t, stdout.Bytes())
	require.NotEmpty(t, stderr.Bytes())
}
