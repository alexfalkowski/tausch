package exec

import (
	"context"
	"os"
	"os/exec"
)

// Command returns an [exec.Cmd] that runs the given command through the
// `tausch` CLI.
//
// It mirrors [os/exec.Command] for callers that do not need context
// cancellation. The returned command executes the `tausch` binary and prefixes
// the provided command with `--`, so this:
//
//	exec.Command("go", "version")
//
// results in an invocation equivalent to:
//
//	tausch -- go version
//
// Tausch CLI flags such as `-config` must be supplied through `TAUSCH_CONFIG` or
// the default config location instead of through name or arg.
func Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(executable(), commandArgs(name, arg...)...) //nolint:noctx // Mirrors os/exec.Command.
}

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
// The provided name and args are target command tokens only. They are always
// passed after `--`, so tausch CLI flags such as `-config` must be supplied
// through `TAUSCH_CONFIG` or the default config location instead.
//
// Executable resolution follows the tausch CLI wrapper rules:
//   - Prefer a `tausch` binary found on PATH.
//   - Otherwise fall back to the `TAUSCH_PATH` environment variable.
//
// PATH lookup requires os/exec.LookPath to return a usable executable path. If
// Go rejects a relative current-directory match such as PATH=".", this function
// falls back to `TAUSCH_PATH`.
//
// If neither provides a usable executable, running the returned command will
// fail with an underlying exec error (for example “executable file not found”).
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, executable(), commandArgs(name, arg...)...)
}

func executable() string {
	path, err := exec.LookPath("tausch")
	if err != nil {
		return os.Getenv("TAUSCH_PATH")
	}

	return path
}

func commandArgs(name string, arg ...string) []string {
	return append([]string{"--", name}, arg...)
}
