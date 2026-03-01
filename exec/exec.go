package exec

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

// CommandContext returns an [exec.Cmd] that runs the given command through the
// `tausch` CLI.
//
// The returned command executes the `tausch` binary and prefixes the provided
// command with `--`, so this:
//
//	exec.CommandContext(ctx, "go", "version")
//
// results in an invocation equivalent to:
//
//	tausch -- go version
//
// This allows command execution to be stubbed by tausch according to its YAML
// configuration.
//
// Executable resolution follows the tausch CLI wrapper rules:
//   - Prefer a `tausch` binary found on PATH.
//   - Otherwise fall back to the `TAUSCH_PATH` environment variable.
//
// If neither provides a usable executable, running the returned command will
// fail with an underlying exec error (for example “executable file not found”).
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, executable(), args(name, arg...)...)
}

func executable() string {
	path, err := exec.LookPath("tausch")
	if err != nil && !errors.Is(err, exec.ErrDot) {
		return os.Getenv("TAUSCH_PATH")
	}

	return path
}

func args(name string, arg ...string) []string {
	return append([]string{"--", name}, arg...)
}
