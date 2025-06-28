package exec_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/exec"
	"github.com/stretchr/testify/require"
)

func init() {
	exec.Register("../tausch", "../test/configs/exec.yaml")
}

func TestCommandSuccess(t *testing.T) {
	b := &bytes.Buffer{}

	cmd := exec.Command("go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandError(t *testing.T) {
	b := &bytes.Buffer{}

	cmd := exec.Command("go", "bob")
	cmd.Stdout = b
	cmd.Stderr = b

	require.Error(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandContextSuccess(t *testing.T) {
	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandContextError(t *testing.T) {
	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "bob")
	cmd.Stdout = b
	cmd.Stderr = b

	require.Error(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}
