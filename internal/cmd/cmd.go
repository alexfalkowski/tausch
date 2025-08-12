package cmd

import (
	"fmt"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/alexfalkowski/tausch/internal/io"
)

// Run will get the command and write to the specified writer.
func Run(stdout, stderr io.Writer, args []string) (int, error) {
	f, err := flag.NewValues(args)
	if err != nil {
		return 0, err
	}

	file, err := f.Config()
	if err != nil {
		return 0, err
	}

	config, err := config.Decode(file)
	if err != nil {
		return 0, err
	}

	name := f.Name()
	command, err := config.GetCommand(name)
	if err != nil {
		return 0, fmt.Errorf("find %s: %w", name, err)
	}

	ok, err := io.Write(stdout, command.Stdout)
	if err != nil || ok {
		return 0, err
	}

	_, err = io.Write(stderr, command.Stderr)
	return 1, err
}
