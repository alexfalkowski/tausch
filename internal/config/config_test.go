package config_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDecodeSuccess(t *testing.T) {
	c, err := config.Decode("../../test/configs/config.yaml")

	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestDecodeError(t *testing.T) {
	values := []string{
		"../../test/configs/none.yaml",
		"../../test/configs/empty.yaml",
	}

	for _, value := range values {
		c, err := config.Decode(value)

		assert.Error(t, err)
		assert.Nil(t, c)
	}
}
