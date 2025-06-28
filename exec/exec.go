package exec

import (
	"context"
	"os"
	"os/exec"
)

var (
	path   = os.Getenv("TAUSCH_PATH")
	config = os.Getenv("TAUSCH_CONFIG")
)

// Register by passing the paths (good for testing).
func Register(p, c string) {
	path = p
	config = c
}

// Command will use tausch under the hood so it looks like you are using the same command.
func Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(path, args(name, arg...)...)
}

// CommandContext is same as Command with a context.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, path, args(name, arg...)...)
}

func args(name string, arg ...string) []string {
	args := []string{"-config", config, "--"}
	args = append(args, name)
	args = append(args, arg...)

	return args
}
