package exec_test

import (
	"context"
	"fmt"
	"os"

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
	// go version go1.24.4 darwin/amd64
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
	// go version go1.24.4 darwin/amd64
}

func setupExampleTausch() (func(), error) {
	dir, err := os.MkdirTemp("", "tausch-example")
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	oldPath := os.Getenv("PATH")
	oldTausch, hadTausch := os.LookupEnv("TAUSCH_PATH")
	oldConfig, hadConfig := os.LookupEnv("TAUSCH_CONFIG")

	os.Setenv("PATH", dir)
	os.Setenv("TAUSCH_PATH", "../tausch")
	os.Setenv("TAUSCH_CONFIG", "../test/configs/config.yml")

	return func() {
		os.Setenv("PATH", oldPath)
		restoreEnv("TAUSCH_PATH", oldTausch, hadTausch)
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
