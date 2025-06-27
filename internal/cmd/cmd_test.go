package cmd_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/tausch/internal/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRunMissingConfig(t *testing.T) {
	os.Args = []string{"cmd", "-config", "cfg.json", "test", "my", "code"}
	c, err := cmd.Run()

	assert.Error(t, err)
	assert.Zero(t, c)
}

func TestRunMissingCommand(t *testing.T) {
	os.Args = []string{"cmd", "-config", "../../test/configs/config.json", "test", "my", "code"}
	c, err := cmd.Run()

	assert.Error(t, err)
	assert.Zero(t, c)
}

func TestRunStdout(t *testing.T) {
	os.Args = []string{"cmd", "-config", "../../test/configs/config.json", "text_stdout"}
	c, err := cmd.Run()

	assert.NoError(t, err)
	assert.Zero(t, c)
}

func TestRunStderr(t *testing.T) {
	os.Args = []string{"cmd", "-config", "../../test/configs/config.json", "text_stderr"}
	c, err := cmd.Run()

	assert.NoError(t, err)
	assert.Equal(t, 1, c)
}
