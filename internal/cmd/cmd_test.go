package cmd_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRunInvalidArgs(t *testing.T) {
	args := []string{"- x"}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	assert.Error(t, err)
	assert.Zero(t, c)
	assert.Empty(t, b.Bytes())
}

func TestRunMissingConfig(t *testing.T) {
	args := []string{
		"-config", "cfg.yaml",
		"--",
		"test", "my", "code",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	assert.Error(t, err)
	assert.Zero(t, c)
	assert.Empty(t, b.Bytes())
}

func TestRunMissingCommand(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yaml",
		"--",
		"test", "my", "code",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	assert.Error(t, err)
	assert.Zero(t, c)
	assert.Empty(t, b.Bytes())
}

func TestRunStdout(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yaml",
		"--",
		"go", "version",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	assert.NoError(t, err)
	assert.Zero(t, c)
	assert.NotEmpty(t, b.Bytes())
}

func TestRunStderr(t *testing.T) {
	args := []string{
		"-config", "../../test/configs/config.yaml",
		"--",
		"go", "bob",
	}
	b := &bytes.Buffer{}
	c, err := cmd.Run(b, b, args)

	assert.NoError(t, err)
	assert.Equal(t, 1, c)
	assert.NotEmpty(t, b.Bytes())
}
