package flag

import (
	"cmp"
	"flag"
	"os"
	"path"
	"strings"
)

var file string

func init() {
	flag.CommandLine.Init("tausch", flag.ContinueOnError)
	flag.CommandLine.StringVar(&file, "config", "", "the config file path")
}

// Config file.
func Config(args []string) (string, error) {
	if err := flag.CommandLine.Parse(args); err != nil {
		return "", err
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	config := cmp.Or(
		file,
		os.Getenv("TAUSCH_CONFIG"),
		path.Join(dir, "tausch", "config.yml"),
	)
	return config, nil
}

// Name of the command.
func Name() string {
	return strings.TrimSpace(strings.Join(flag.CommandLine.Args(), " "))
}
