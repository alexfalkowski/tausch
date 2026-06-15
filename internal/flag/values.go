package flag

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// Parse parses args and returns parsed [Values].
//
// Parsing uses a dedicated FlagSet configured with flag.ContinueOnError, so
// parsing errors are returned to the caller rather than terminating the process.
// Any parse error or usage output produced by the FlagSet is written to output.
//
// The tausch CLI is typically invoked as:
//
//	tausch -config path/to/config.yml -- <command tokens...>
//
// After parsing flags, the remaining arguments (typically those after `--`) are
// preserved and used to derive the command name via [Values.Name].
func Parse(output io.Writer, args []string) (*Values, error) {
	set := flag.NewFlagSet("tausch", flag.ContinueOnError)
	set.SetOutput(output)
	set.Usage = func() {
		fmt.Fprintln(set.Output(), "Usage of tausch:")
		fmt.Fprintln(set.Output(), "  tausch [flags] -- <command tokens...>")
		fmt.Fprintln(set.Output())
		fmt.Fprintln(set.Output(), "Command tokens after -- are joined with spaces")
		fmt.Fprintln(set.Output(), "and matched against a config entry's name.")
		fmt.Fprintln(set.Output())
		fmt.Fprintln(set.Output(), "Config path resolution:")
		fmt.Fprintln(set.Output(), "  1. -config <path>")
		fmt.Fprintln(set.Output(), "  2. TAUSCH_CONFIG")
		fmt.Fprintln(set.Output(), "  3. os.UserConfigDir()/tausch/config.yml")
		fmt.Fprintln(set.Output())
		fmt.Fprintln(set.Output(), "Flags:")
		set.PrintDefaults()
	}

	file := set.String("config", "", "the config file path")
	if err := set.Parse(args); err != nil {
		return nil, err
	}

	return &Values{file: *file, args: set.Args()}, nil
}

// Values represents the parsed CLI inputs relevant to tausch.
//
// It stores the resolved flag value for the config path (if provided) and the
// remaining arguments used to derive the command name for config lookup.
type Values struct {
	file string
	args []string
}

// Config resolves the tausch YAML config path based on precedence rules.
//
// Resolution order:
//
//  1. The `-config` flag (if provided).
//  2. The `TAUSCH_CONFIG` environment variable (if set).
//  3. A default under the user config directory:
//     $HOME/.config/tausch/config.yml (or platform equivalent from os.UserConfigDir).
//
// This method selects a path but does not validate that it exists or is readable;
// callers should expect failures when opening/decoding the file if the path is
// invalid.
func (f *Values) Config() (string, error) {
	if f.file != "" {
		return f.file, nil
	}

	if config := os.Getenv("TAUSCH_CONFIG"); config != "" {
		return config, nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "tausch", "config.yml"), nil
}

// Name returns the command name derived from the remaining (non-flag) arguments.
//
// The name is produced by joining the remaining args with a single space and
// trimming leading/trailing whitespace.
//
// For example, for:
//
//	tausch -- go version
//
// the derived name is:
//
//	"go version"
//
// This name must match the YAML config command entry `name` field exactly.
func (f *Values) Name() string {
	return strings.TrimSpace(strings.Join(f.args, " "))
}
