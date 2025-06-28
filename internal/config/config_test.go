package config_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/stretchr/testify/require"
)

func TestDecodeSuccess(t *testing.T) {
	c, err := config.Decode("../../test/configs/config.yaml")

	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestDecodeError(t *testing.T) {
	values := []string{
		"../../test/configs/none.yaml",
		"../../test/configs/empty.yaml",
	}

	for _, value := range values {
		c, err := config.Decode(value)

		require.Error(t, err)
		require.Nil(t, c)
	}
}
