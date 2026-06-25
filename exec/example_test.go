package exec_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/tausch/exec"
)

func ExampleCommand() {
	cleanup, err := setupExampleTausch()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanup()

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(string(out))
	// Output:
	// -- go version
}

func ExampleCommandContext() {
	cleanup, err := setupExampleTausch()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanup()

	out, err := exec.CommandContext(context.Background(), "go", "version").Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(string(out))
	// Output:
	// -- go version
}

func setupExampleTausch() (func(), error) {
	dir, err := os.MkdirTemp("", "tausch-example")
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	tausch := filepath.Join(dir, "tausch")
	if err := os.WriteFile(tausch, []byte("#!/bin/sh\nprintf '%s\\n' \"$*\"\n"), 0o600); err != nil {
		cleanup()
		return nil, err
	}
	if err := os.Chmod(tausch, 0o700); err != nil {
		cleanup()
		return nil, err
	}

	oldPath := os.Getenv("PATH")
	oldConfig, hadConfig := os.LookupEnv("TAUSCH_CONFIG")
	os.Setenv("PATH", dir)
	os.Setenv("TAUSCH_CONFIG", "config.yml")

	return func() {
		os.Setenv("PATH", oldPath)
		restoreEnv("TAUSCH_CONFIG", oldConfig, hadConfig)
		cleanup()
	}, nil
}

func restoreEnv(key, value string, hadValue bool) {
	if hadValue {
		os.Setenv(key, value)
		return
	}

	os.Unsetenv(key)
}
