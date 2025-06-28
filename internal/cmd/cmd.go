package cmd

import (
	"fmt"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/alexfalkowski/tausch/internal/io"
)

// Run will get the command and write to the specified writer.
func Run(stdout, stderr io.Writer) (int, error) {
	file := flag.Config()
	name := flag.Name()

	config, err := config.Decode(file)
	if err != nil {
		return 0, err
	}

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
