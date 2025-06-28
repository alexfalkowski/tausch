package exec_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/alexfalkowski/tausch/exec"
	"github.com/stretchr/testify/require"
)

func TestPathCommandSuccess(t *testing.T) {
	path := os.Getenv("PATH") + ":../"

	t.Setenv("PATH", path)
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	b := &bytes.Buffer{}

	cmd := exec.Command("go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandSuccess(t *testing.T) {
	t.Setenv("TAUSCH_PATH", "../tausch")
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	b := &bytes.Buffer{}

	cmd := exec.Command("go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandError(t *testing.T) {
	t.Setenv("TAUSCH_PATH", "../tausch")
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	b := &bytes.Buffer{}

	cmd := exec.Command("go", "bob")
	cmd.Stdout = b
	cmd.Stderr = b

	require.Error(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandContextSuccess(t *testing.T) {
	t.Setenv("TAUSCH_PATH", "../tausch")
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandContextError(t *testing.T) {
	t.Setenv("TAUSCH_PATH", "../tausch")
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "bob")
	cmd.Stdout = b
	cmd.Stderr = b

	require.Error(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}
