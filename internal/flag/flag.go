package flag

import (
	"flag"
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

	return file, nil
}

// Name of the command.
func Name() string {
	return strings.TrimSpace(strings.Join(flagSet.Args(), " "))
}
