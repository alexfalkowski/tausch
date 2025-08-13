package flag

import (
	"cmp"
	"flag"
	"os"
	"path"
	"strings"
)

// Parse parses the args and extracts the config file and command name.
func Parse(args []string) (*Values, error) {
	set := flag.NewFlagSet("tausch", flag.ContinueOnError)

	file := set.String("config", "", "the config file path")
	if err := set.Parse(args); err != nil {
		return nil, err
	}

	return &Values{file: *file, args: set.Args()}, nil
}

// Values represents the config and name of the command.
type Values struct {
	file string
	args []string
}

// Config file.
func (f *Values) Config() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	config := cmp.Or(
		f.file,
		os.Getenv("TAUSCH_CONFIG"),
		path.Join(dir, "tausch", "config.yml"),
	)
	return config, nil
}

// Name of the command.
func (f *Values) Name() string {
	return strings.TrimSpace(strings.Join(f.args, " "))
}
