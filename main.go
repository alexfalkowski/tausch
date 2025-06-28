package main

import (
	"log"
	"os"

	"github.com/alexfalkowski/tausch/internal/cmd"
)

func main() {
	code, err := cmd.Run(os.Stdout, os.Stderr, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
