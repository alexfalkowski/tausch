package flag_test

import (
	"io"
	"os"
	"path"
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

func TestValuesConfig(t *testing.T) {
	tests := []configCase{
		{
			name:      "flag takes precedence with missing user config dir",
			envConfig: "env.yml",
			want:      "cfg.yml",
			args:      []string{"-config", "cfg.yml"},
		},
		{
			name:      "env takes precedence with missing user config dir",
			envConfig: "env.yml",
			want:      "env.yml",
		},
		{
			name:         "default errors with missing user config dir",
			requireError: true,
		},
		{
			name:        "default uses user config dir",
			useHome:     true,
			wantDefault: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t)
		})
	}
}

type configCase struct {
	name         string
	envConfig    string
	want         string
	args         []string
	useHome      bool
	wantDefault  bool
	requireError bool
}

func (c configCase) assert(t *testing.T) {
	t.Helper()

	if c.useHome {
		t.Setenv("HOME", t.TempDir())
	} else {
		t.Setenv("HOME", "")
	}
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("TAUSCH_CONFIG", c.envConfig)

	f, err := flag.Parse(io.Discard, c.args)
	require.NoError(t, err)

	config, err := f.Config()
	if c.requireError {
		require.Error(t, err)
		require.Empty(t, config)
		return
	}

	require.NoError(t, err)
	if c.wantDefault {
		configDir, err := os.UserConfigDir()
		require.NoError(t, err)
		require.Equal(t, path.Join(configDir, "tausch", "config.yml"), config)
		return
	}

	require.Equal(t, c.want, config)
}

func TestValuesError(t *testing.T) {
	args := []string{"- x"}

	_, err := flag.Parse(io.Discard, args)
	require.Error(t, err)
}
