package exec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

// Command will use tausch under the hood so it looks like you are using the same command.
func Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(executable(), args(name, arg...)...)
}

// CommandContext is same as Command with a context.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, executable(), args(name, arg...)...)
}

func executable() string {
	execPath := os.Getenv("TAUSCH_PATH")

	path, err := exec.LookPath("tausch")
	if err != nil {
		return execPath
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return execPath
	}

	return path
}

func args(name string, arg ...string) []string {
	args := []string{"--", name}
	args = append(args, arg...)

	return args
}
