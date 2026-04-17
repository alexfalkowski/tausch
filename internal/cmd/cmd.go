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
//   - If argument parsing fails, the flag package writes usage/error output to
//     stderr and Run returns 1.
//   - If config resolution, config decoding, command lookup, or stdout/stderr
//     decode/write steps fail, Run writes the error to stderr and returns 1.
//   - If the command's configured `stdout` is non-empty and is successfully
//     written, Run returns 0.
//   - Otherwise, Run writes the configured `stderr` payload and returns 1.
//
// Run owns user-facing error output for the CLI flow; callers should treat the
// returned status code as the complete result.
func Run(stdout, stderr io.Writer, args []string) int {
	f, err := flag.Parse(stderr, args)
	if err != nil {
		return 1
	}

	file, err := f.Config()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	config, err := config.Decode(file)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	name := f.Name()
	command, err := config.GetCommand(name)
	if err != nil {
		fmt.Fprintln(stderr, fmt.Errorf("find %s: %w", name, err))
		return 1
	}

	ok, err := io.Write(stdout, command.Stdout)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if ok {
		return 0
	}

	_, err = io.Write(stderr, command.Stderr)
	if err != nil {
		fmt.Fprintln(stderr, err)
	}

	return 1
}
