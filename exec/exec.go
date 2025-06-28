package exec

import (
	"context"
	"fmt"
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

// CommandContext will use tausch under the hood so it looks like you are using the same command.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	args := []string{"-config", config, "--"}
	args = append(args, name)
	args = append(args, arg...)

	fmt.Println(args)

	return exec.CommandContext(ctx, path, args...)
}
