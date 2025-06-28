package cmd_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/cmd"
	"github.com/stretchr/testify/require"
)

func TestRunInvalidArgs(t *testing.T) {
	args := []string{"- x"}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	require.Error(t, err)
	require.Zero(t, c)
	require.Empty(t, b.Bytes())
}

func TestRunMissingConfig(t *testing.T) {
	args := []string{
		"-config", "cfg.yml",
		"--",
		"test", "my", "code",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	require.Error(t, err)
	require.Zero(t, c)
	require.Empty(t, b.Bytes())
}

func TestRunMissingCommand(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"test", "my", "code",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	require.Error(t, err)
	require.Zero(t, c)
	require.Empty(t, b.Bytes())
}

func TestRunStdout(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"go", "version",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	require.NoError(t, err)
	require.Zero(t, c)
	require.NotEmpty(t, b.Bytes())
}

func TestRunStderr(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"go", "bob",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	require.NoError(t, err)
	require.Equal(t, 1, c)
	require.NotEmpty(t, b.Bytes())
}
