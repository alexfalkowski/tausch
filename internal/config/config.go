package config

import (
	"errors"
	"os"
	"slices"

	"go.yaml.in/yaml/v3"
)

// ErrCommandNotFound for config.
var ErrCommandNotFound = errors.New("command not found")

// Decode the config.
func Decode(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Config for tausch.
type Config struct {
	Cmds []*Command `json:"cmds"`
}

// GetCommand by name.
func (c *Config) GetCommand(name string) (*Command, error) {
	i := slices.IndexFunc(c.Cmds, func(c *Command) bool { return c.Name == name })
	if i == -1 {
		return nil, ErrCommandNotFound
	}

	return c.Cmds[i], nil
}

// Command for config.
type Command struct {
	Name   string `json:"name"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
