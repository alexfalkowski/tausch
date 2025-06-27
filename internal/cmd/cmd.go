package cmd

import (
	"fmt"
	"os"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/alexfalkowski/tausch/internal/io"
)

// Run will get the command as write to the correct writer.
func Run() (int, error) {
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

	ok, err := io.Write(command.Stdout, os.Stdout)
	if err != nil || ok {
		return 0, err
	}

	_, err = io.Write(command.Stderr, os.Stderr)

	return 1, err
}
