package cmd_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/tausch/internal/cmd"
	"github.com/alexfalkowski/tausch/internal/test"
	"github.com/stretchr/testify/require"
)

func TestRunInvalidArgs(t *testing.T) {
	t.Parallel()

	args := []string{"- x"}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "flag provided but not defined")
	require.Contains(t, stderr.String(), "Usage of tausch:")
	require.Contains(t, stderr.String(), "tausch [flags] -- <command tokens...>")
	require.Contains(t, stderr.String(), "TAUSCH_CONFIG")
}

func TestRunConfigError(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("TAUSCH_CONFIG", "")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, nil)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "$HOME")
}

func TestRunMissingConfig(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "cfg.yml",
		"--",
		"test", "my", "code",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "cfg.yml")
}

func TestRunMissingCommand(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "../../test/configs/config.yml",
		"--",
		"test", "my", "code",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "find test my code: command not found")
}

func TestRunMultipleOutputs(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "../../test/configs/multiple_outputs.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "multiple outputs configured")
}

func TestRunStdoutWriteError(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "../../test/configs/stdout_invalid_base64.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "illegal base64 data")
}

func TestRunUnselectedMalformedPayload(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "../../test/configs/unselected_invalid_base64.yml",
		"--",
		"demo",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Zero(t, code)
	require.Equal(t, "valid", stdout.String())
	require.Empty(t, stderr.String())
}

func TestRunStdoutWriterError(t *testing.T) {
	t.Chdir("../..")

	args := []string{
		"-config", "test/configs/config.yml",
		"--",
		"go", "version",
	}
	stderr := &bytes.Buffer{}
	code := cmd.Run(test.FailingWriter{}, stderr, args)

	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), test.ErrWriteFailed.Error())
}

func TestRunStdout(t *testing.T) {
	t.Chdir("../..")

	args := []string{
		"-config", "test/configs/config.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Zero(t, code)
	require.NotEmpty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunStdoutExitCode(t *testing.T) {
	t.Chdir("../..")

	args := []string{
		"-config", "test/configs/exit_code.yml",
		"--",
		"go", "version",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 7, code)
	require.NotEmpty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunStderrWriteError(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "../../test/configs/stderr_invalid_base64.yml",
		"--",
		"go", "bob",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), "illegal base64 data")
}

func TestRunStderr(t *testing.T) {
	t.Chdir("../..")

	args := []string{
		"-config", "test/configs/config.yml",
		"--",
		"go", "bob",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.NotEmpty(t, stderr.Bytes())
}

func TestRunStderrExitCode(t *testing.T) {
	t.Chdir("../..")

	args := []string{
		"-config", "test/configs/exit_code.yml",
		"--",
		"go", "bob",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 127, code)
	require.Empty(t, stdout.Bytes())
	require.NotEmpty(t, stderr.Bytes())
}

func TestRunNoOutputExitCode(t *testing.T) {
	t.Parallel()

	args := []string{
		"-config", "../../test/configs/exit_code.yml",
		"--",
		"grep", "needle", "file.txt",
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cmd.Run(stdout, stderr, args)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())
}

func TestRunStatusEdges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		stdout     string
		stderrPart string
		args       []string
		code       int
	}{
		{
			name:   "silent stdout",
			args:   []string{"-config", "../../test/configs/status_edges.yml", "--", "silent stdout"},
			code:   0,
			stdout: "",
		},
		{
			name:   "no output success",
			args:   []string{"-config", "../../test/configs/status_edges.yml", "--", "no output success"},
			code:   0,
			stdout: "",
		},
		{
			name:   "no output maximum",
			args:   []string{"-config", "../../test/configs/status_edges.yml", "--", "no output maximum"},
			code:   255,
			stdout: "",
		},
		{
			name:       "invalid output overrides exit code",
			args:       []string{"-config", "../../test/configs/status_edges.yml", "--", "invalid output"},
			code:       1,
			stdout:     "",
			stderrPart: "illegal base64 data",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			code := cmd.Run(stdout, stderr, test.args)

			require.Equal(t, test.code, code)
			require.Equal(t, test.stdout, stdout.String())
			if test.stderrPart == "" {
				require.Empty(t, stderr.String())
			} else {
				require.Contains(t, stderr.String(), test.stderrPart)
			}
		})
	}
}
