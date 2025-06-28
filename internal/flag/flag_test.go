package flag_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/stretchr/testify/require"
)

func TestConfigSuccess(t *testing.T) {
	args := []string{"-config", "cfg.yml", "test", "my", "code"}

	c, err := flag.Config(args)
	require.Equal(t, "cfg.yml", c)
	require.NoError(t, err)

	name := flag.Name()
	require.Equal(t, "test my code", name)
}

func TestConfigError(t *testing.T) {
	args := []string{"- x"}

	c, err := flag.Config(args)
	require.Empty(t, c)
	require.Error(t, err)

	name := flag.Name()
	require.Empty(t, name)
}
