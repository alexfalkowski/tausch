package flag_test

import (
	"testing"

	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/stretchr/testify/assert"
)

func TestConfigSuccess(t *testing.T) {
	args := []string{"-config", "cfg.yaml", "test", "my", "code"}

	c, err := flag.Config(args)
	assert.Equal(t, "cfg.yaml", c)
	assert.NoError(t, err)

	name := flag.Name()
	assert.Equal(t, "test my code", name)
}

func TestConfigError(t *testing.T) {
	args := []string{"- x"}

	c, err := flag.Config(args)
	assert.Empty(t, c)
	assert.Error(t, err)

	name := flag.Name()
	assert.Empty(t, name)
}
