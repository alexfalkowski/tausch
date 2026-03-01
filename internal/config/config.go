package config

import (
	"errors"
	"os"
	"slices"

	"go.yaml.in/yaml/v3"
)

// ErrCommandNotFound is returned when a command lookup fails because no entry in
// the decoded configuration has a matching [Command.Name].
//
// Callers typically compare with errors.Is(err, ErrCommandNotFound).
var ErrCommandNotFound = errors.New("command not found")

// Decode reads a YAML configuration file from path and decodes it into a [Config].
//
// The YAML is expected to match the tausch schema (a top-level `cmds` list). This
// function only performs YAML decoding; it does not validate that any commands
// are present or that stdout/stderr payload strings are valid `kind:data` values.
//
// The returned *Config is ready to be queried with [Config.GetCommand].
//
// Errors are returned for I/O failures (opening/reading the file) and for YAML
// decoding errors.
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

	return &config, nil
}

// Config is the in-memory representation of the tausch YAML configuration.
//
// It corresponds to a YAML document with a top-level `cmds` key containing a list
// of command stubs.
//
// The stdout/stderr fields inside each command are stored as opaque strings here;
// decoding those values into bytes is handled by internal/encoding and internal/io.
type Config struct {
	// Cmds is the set of configured command stubs.
	Cmds []*Command `yaml:"cmds"`
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
	i := slices.IndexFunc(c.Cmds, func(c *Command) bool { return c.Name == name })
	if i == -1 {
		return nil, ErrCommandNotFound
	}

	return c.Cmds[i], nil
}

// Command is a single command stub entry in the configuration.
//
// Name identifies the command to match (as derived by the CLI), and Stdout/Stderr
// define what tausch should emit when that command is invoked.
//
// Stdout and Stderr are strings using tausch's `kind:data` format. Supported kinds
// are decoded by internal/encoding (for example `text:...`, `file:...`, `base64:...`).
type Command struct {
	// Name is the exact command name to match (for example "go version").
	Name string `yaml:"name"`

	// Stdout is the encoded payload to write to stdout when the command is treated
	// as successful (non-empty value).
	Stdout string `yaml:"stdout"`

	// Stderr is the encoded payload to write to stderr when stdout is empty and the
	// command is treated as failing.
	Stderr string `yaml:"stderr"`
}
