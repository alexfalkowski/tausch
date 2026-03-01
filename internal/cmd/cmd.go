package cmd

import (
	"fmt"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/alexfalkowski/tausch/internal/io"
)

// Run executes the tausch CLI flow for the provided arguments and writes the
// configured output to the supplied writers.
//
// It performs the following steps:
//
//  1. Parses flags and remaining arguments from args (via internal/flag).
//  2. Resolves the YAML config path (via (*flag.Values).Config).
//  3. Decodes the YAML configuration (via internal/config.Decode).
//  4. Derives the command name (via (*flag.Values).Name) and looks up the
//     matching command entry (via (*config.Config).GetCommand).
//  5. Writes the configured `stdout` payload to stdout if present; otherwise it
//     writes the configured `stderr` payload to stderr (via internal/io.Write).
//
// The configured stdout/stderr payloads are strings in tausch's `kind:data`
// format (for example `text:...`, `file:...`, `base64:...`) and are decoded by
// internal/io before being written.
//
// # Exit code contract
//
// Run returns an integer exit code intended to be used as the process exit code:
//
//   - If the command's configured `stdout` is non-empty and is successfully
//     written, Run returns (0, nil).
//   - Otherwise, Run attempts to write the configured `stderr` and returns
//     (1, err) where err is any error from decoding or writing stderr.
//
// Any error encountered while parsing flags, resolving the config path, decoding
// YAML, or looking up the command is returned with an exit code of 0 (callers
// should treat a non-nil error as failure regardless of the code).
func Run(stdout, stderr io.Writer, args []string) (int, error) {
	f, err := flag.Parse(args)
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
