package flag_test

import (
	"io"
	"testing"

	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/stretchr/testify/require"
)

func TestValuesSuccess(t *testing.T) {
	args := []string{"-config", "cfg.yml", "test", "my", "code"}

	f, err := flag.Parse(io.Discard, args)
	require.NoError(t, err)

	c, err := f.Config()
	require.Equal(t, "cfg.yml", c)
	require.NoError(t, err)

	name := f.Name()
	require.Equal(t, "test my code", name)
}

func TestValuesError(t *testing.T) {
	args := []string{"- x"}

	_, err := flag.Parse(io.Discard, args)
	require.Error(t, err)
}
