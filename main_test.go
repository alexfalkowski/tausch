package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadmeCLIExamples(t *testing.T) {
	tests := []struct {
		name       string
		wantStdout string
		wantStderr string
		args       []string
		wantCode   int
	}{
		{
			name:       "stdout",
			wantStdout: "test/stdout/go_version.txt",
			wantCode:   0,
			args: []string{
				"-config", "test/configs/config.yml",
				"--",
				"go", "version",
			},
		},
		{
			name:       "stderr",
			wantStderr: "test/stderr/go_bob.txt",
			wantCode:   1,
			args: []string{
				"-config", "test/configs/config.yml",
				"--",
				"go", "bob",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cmd := exec.CommandContext(t.Context(), "./tausch", tt.args...)
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			err := cmd.Run()
			if tt.wantCode == 0 {
				require.NoError(t, err)
			} else {
				var exitErr *exec.ExitError
				require.ErrorAs(t, err, &exitErr)
				require.Equal(t, tt.wantCode, exitErr.ExitCode())
			}

			require.Equal(t, readFixture(t, tt.wantStdout), stdout.Bytes())
			require.Equal(t, readFixture(t, tt.wantStderr), stderr.Bytes())
		})
	}
}

func readFixture(t *testing.T, path string) []byte {
	t.Helper()

	if path == "" {
		return []byte{}
	}

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return data
}
