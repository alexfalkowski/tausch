package exec

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

// CommandContext returns the [Cmd] struct to execute the named program with the given arguments.
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
