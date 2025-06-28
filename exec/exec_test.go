package exec_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/exec"
	"github.com/stretchr/testify/assert"
)

func init() {
	exec.Register("../tausch", "../test/configs/exec.yaml")
}

func TestCommandContextSuccess(t *testing.T) {
	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	assert.NoError(t, cmd.Run())
	assert.NoError(t, cmd.Err)
	assert.NotEmpty(t, b.Bytes())
}

func TestCommandContextError(t *testing.T) {
	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "bob")
	cmd.Stdout = b
	cmd.Stderr = b

	assert.Error(t, cmd.Run())
	assert.NoError(t, cmd.Err)
	assert.NotEmpty(t, b.Bytes())
}
