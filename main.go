package main

import (
	"os"

	"github.com/alexfalkowski/tausch/internal/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Stdout, os.Stderr, os.Args[1:]))
}
