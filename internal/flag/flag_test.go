package flag_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Args = []string{"cmd", "-config", "cfg.yaml", "test", "my", "code"}

	c := flag.Config()
	assert.Equal(t, "cfg.yaml", c)

	name := flag.Name()
	assert.Equal(t, "test my code", name)
}
