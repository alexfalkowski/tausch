package flag

import (
	"cmp"
	"flag"
	"os"
	"path"
	"strings"
)

var (
	flagSet *flag.FlagSet
	file    string
)

func init() {
	flagSet = flag.NewFlagSet("tausch", flag.ContinueOnError)
	flagSet.StringVar(&file, "config", "", "the config file path")
}

// Config file.
func Config(args []string) (string, error) {
	if err := flagSet.Parse(args); err != nil {
		return "", err
	}

	home, _ := os.UserHomeDir()
	config := cmp.Or(
		file,
		os.Getenv("TAUSCH_CONFIG"),
		path.Join(home, ".config", "tausch.yml"),
	)

	return config, nil
}

// Name of the command.
func Name() string {
	return strings.TrimSpace(strings.Join(flagSet.Args(), " "))
}
