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

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.Equal(t, readFixture(t, "../test/stdout/go_version.txt"), stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestCommandPrefersPathOverTauschPath(t *testing.T) {
	pathDir := t.TempDir()
	path := filepath.Join(pathDir, "tausch")
	fallback := filepath.Join(t.TempDir(), "tausch")

	writeExecutable(t, path, "#!/bin/sh\nprintf path\n")
	writeExecutable(t, fallback, "#!/bin/sh\nprintf fallback\n")

	t.Setenv("PATH", pathDir)
	t.Setenv("TAUSCH_PATH", fallback)

	b := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = b
	cmd.Stderr = b

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.Equal(t, "path", b.String())
}

func TestCommandSuccess(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("TAUSCH_PATH", "../tausch")
	t.Setenv("TAUSCH_CONFIG", "../test/configs/exec.yml")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "version")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	require.NoError(t, cmd.Run())
	require.NoError(t, cmd.Err)
	require.Equal(t, readFixture(t, "../test/stdout/go_version.txt"), stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestCommandFallsBackToTauschPathForErrDot(t *testing.T) {
	cwd := t.TempDir()
	fallback := filepath.Join(t.TempDir(), "tausch")
	local := filepath.Join(cwd, "tausch")

	writeExecutable(t, fallback, "#!/bin/sh\nprintf fallback\n")
	writeExecutable(t, local, "#!/bin/sh\nprintf local\n")

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

func TestCommandPrefixesDelimiter(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("TAUSCH_PATH", filepath.Join(t.TempDir(), "tausch"))

	cmd := exec.CommandContext(t.Context(), "-version", "--json")

	require.Equal(t, []string{"--", "-version", "--json"}, cmd.Args[1:])
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

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.CommandContext(t.Context(), "go", "bob")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	require.Error(t, err)
	require.NoError(t, cmd.Err)
	require.NotNil(t, cmd.ProcessState)
	require.Equal(t, 1, cmd.ProcessState.ExitCode())
	require.Empty(t, stdout.Bytes())
	require.Equal(t, readFixture(t, "../test/stderr/go_bob.txt"), stderr.Bytes())
}

func readFixture(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return data
}

func writeExecutable(t *testing.T, path, script string) {
	t.Helper()

	require.NoError(t, os.WriteFile(path, []byte(script), 0o600))
	require.NoError(t, os.Chmod(path, 0o700))
}
