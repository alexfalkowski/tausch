package main

import (
	"log"
	"os"

	"github.com/alexfalkowski/tausch/internal/cmd"
)

func main() {
	code, err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
