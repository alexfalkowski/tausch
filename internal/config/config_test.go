package config_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/stretchr/testify/require"
)

func TestDecodeSuccess(t *testing.T) {
	c, err := config.Decode("../../test/configs/config.yml")

	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestDecodeError(t *testing.T) {
	values := []string{
		"../../test/configs/none.yml",
		"../../test/configs/empty.yml",
	}

	for _, value := range values {
		c, err := config.Decode(value)

		require.Error(t, err)
		require.Nil(t, c)
	}
}

func TestDecodeMultipleOutputs(t *testing.T) {
	c, err := config.Decode("../../test/configs/multiple_outputs.yml")

	require.Nil(t, c)
	require.ErrorIs(t, err, config.ErrMultipleOutputs)
	require.Contains(t, err.Error(), `command "go version"`)
}

func TestGetCommandNilCommand(t *testing.T) {
	c := &config.Config{Cmds: []*config.Command{nil}}

	command, err := c.GetCommand("go version")

	require.Nil(t, command)
	require.ErrorIs(t, err, config.ErrCommandNotFound)
}

func TestGetCommandExactMatch(t *testing.T) {
	command := &config.Command{Name: "go version", Stdout: "text:go version"}
	c := &config.Config{
		Cmds: []*config.Command{
			command,
			{Name: "Go Version", Stdout: "text:case"},
			{Name: "go version extra", Stdout: "text:extra"},
		},
	}

	got, err := c.GetCommand("go version")
	require.NoError(t, err)
	require.Same(t, command, got)

	got, err = c.GetCommand("go version extra")
	require.NoError(t, err)
	require.Same(t, c.Cmds[2], got)

	for _, value := range []string{"Go version", "go", " go version "} {
		t.Run(value, func(t *testing.T) {
			got, err := c.GetCommand(value)

			require.Nil(t, got)
			require.ErrorIs(t, err, config.ErrCommandNotFound)
		})
	}
}
