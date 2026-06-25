package config

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"go.yaml.in/yaml/v3"
)

// ErrCommandNotFound is returned when a command lookup fails because no entry in
// the decoded configuration has a matching [Command.Name].
//
// Callers typically compare with errors.Is(err, ErrCommandNotFound).
var ErrCommandNotFound = errors.New("command not found")

// ErrMultipleOutputs is returned when a command configures both stdout and
// stderr. A command stub can emit one configured output stream.
//
// Callers typically compare with errors.Is(err, ErrMultipleOutputs).
var ErrMultipleOutputs = errors.New("multiple outputs configured")

// ErrInvalidExitCode is returned when a command configures an exit code outside
// the supported process exit status range.
//
// Callers typically compare with errors.Is(err, ErrInvalidExitCode).
var ErrInvalidExitCode = errors.New("invalid exit code")

const maxExitCode = 255

// Decode reads a YAML configuration file from path and decodes it into a [Config].
//
// The YAML is expected to match the tausch schema (a top-level `cmds` list). This
// function validates command-level invariants, but it does not validate that any
// commands are present or that stdout/stderr payload strings are valid
// `kind:data` values.
//
// The returned *Config is ready to be queried with [Config.GetCommand].
//
// Errors are returned for I/O failures (opening/reading the file), YAML decoding
// errors, and validation failures such as [ErrMultipleOutputs] or
// [ErrInvalidExitCode].
func Decode(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Config is the in-memory representation of the tausch YAML configuration.
//
// It corresponds to a YAML document with a top-level `cmds` key containing a list
// of command stubs.
//
// The stdout/stderr fields inside each command are stored as opaque strings here;
// decoding those values into bytes is handled by internal/io.
type Config struct {
	// Cmds is the set of configured command stubs.
	Cmds []*Command `yaml:"cmds"`
}

// Validate checks command-level config invariants.
//
// It permits nil command entries so lookup remains tolerant of sparse decoded
// data, but non-nil commands must not configure both stdout and stderr and must
// keep exit codes in range.
func (c *Config) Validate() error {
	for _, command := range c.Cmds {
		if command == nil {
			continue
		}

		if command.hasMultipleOutputs() {
			return fmt.Errorf("command %q: %w", command.Name, ErrMultipleOutputs)
		}

		if command.hasInvalidExitCode() {
			return fmt.Errorf("command %q: %w", command.Name, ErrInvalidExitCode)
		}
	}

	return nil
}

// GetCommand returns the configured [Command] whose [Command.Name] exactly matches
// name.
//
// Matching is string-based and exact (case-sensitive). The tausch CLI derives name
// by joining the command tokens (typically the args after `--`) with spaces.
// Because of this, small differences in whitespace or quoting can cause lookups
// to fail.
//
// If no command matches, the error will be [ErrCommandNotFound].
func (c *Config) GetCommand(name string) (*Command, error) {
	i := slices.IndexFunc(c.Cmds, func(command *Command) bool {
		return command != nil && command.Name == name
	})
	if i == -1 {
		return nil, ErrCommandNotFound
	}

	return c.Cmds[i], nil
}

// Command is a single command stub entry in the configuration.
//
// Name identifies the command to match (as derived by the CLI), Stdout/Stderr
// define what tausch should emit when that command is invoked, and ExitCode can
// override the default process exit code.
//
// Stdout and Stderr are strings using tausch's `kind:data` format. Supported kinds
// are decoded by internal/io (for example `text:...`, `file:...`, `base64:...`).
type Command struct {
	// ExitCode is the optional process exit code to return after successfully
	// writing the configured output.
	ExitCode *int `yaml:"exit_code"`

	// Name is the exact command name to match (for example "go version").
	Name string `yaml:"name"`

	// Stdout is the encoded payload to write to stdout.
	Stdout string `yaml:"stdout"`

	// Stderr is the encoded payload to write to stderr when stdout is empty.
	Stderr string `yaml:"stderr"`
}

func (c *Command) hasMultipleOutputs() bool {
	return c.Stdout != "" && c.Stderr != ""
}

func (c *Command) hasInvalidExitCode() bool {
	return c.ExitCode != nil && (*c.ExitCode < 0 || *c.ExitCode > maxExitCode)
}
