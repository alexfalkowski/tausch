package flag

import (
	"flag"
	"strings"
)

var file string

func init() {
	flag.StringVar(&file, "config", "", "the config file path")
}

// Config file.
func Config() string {
	flag.Parse()

	return file
}

// Name of the command.
func Name() string {
	return strings.TrimSpace(strings.Join(flag.Args(), " "))
}
