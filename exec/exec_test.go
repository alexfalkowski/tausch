package exec_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/tausch/exec"
	"github.com/stretchr/testify/require"
)

func TestPathCommandSuccess(t *testing.T) {
	root, err := filepath.Abs("..")
	require.NoError(t, err)

	t.Setenv("PATH", root)
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.NotEmpty(t, b.Bytes())
}

func TestCommandSuccess(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
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

func TestCommandFallsBackToTauschPathForErrDot(t *testing.T) {
	cwd := t.TempDir()
	fallback := filepath.Join(t.TempDir(), "tausch")
	local := filepath.Join(cwd, "tausch")

	require.NoError(t, os.WriteFile(fallback, []byte("#!/bin/sh\nprintf fallback\n"), 0o600))
	require.NoError(t, os.Chmod(fallback, 0o700))
	require.NoError(t, os.WriteFile(local, []byte("#!/bin/sh\nprintf local\n"), 0o600))
	require.NoError(t, os.Chmod(local, 0o700))

	t.Chdir(cwd)
	t.Setenv("PATH", ".")
	t.Setenv("TAUSCH_PATH", fallback)

	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.Equal(t, "fallback", b.String())
}

func TestCommandPassesVariadicArgs(t *testing.T) {
	config := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(config, []byte(`cmds:
- name: go env GOPATH
  stdout: text:gopath
`), 0o600))

	t.Setenv("PATH", t.TempDir())
	t.Setenv("TAUSCH_PATH", "../tausch")
	t.Setenv("TAUSCH_CONFIG", config)

	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "env", "GOPATH")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.Equal(t, "gopath", b.String())
}

func TestCommandError(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
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
