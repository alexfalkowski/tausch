package exec_test

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/alexfalkowski/tausch/exec"
	"github.com/stretchr/testify/require"
)

func TestHomeCommandSuccess(t *testing.T) {
	t.Setenv("TAUSCH_PATH", "../tausch")

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	dir := path.Join(home, ".config")

	err = os.MkdirAll(dir, 0o600)
	require.NoError(t, err)

	d, err := os.ReadFile("../test/configs/exec.yml")
	require.NoError(t, err)

	config := path.Join(dir, "tausch.yml")

	err = os.WriteFile(config, d, 0o600)
	require.NoError(t, err)

	b := &bytes.Buffer{}

	cmd := exec.Command("go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())

	err = os.Remove(config)
	require.NoError(t, err)
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
