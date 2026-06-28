package flag_test

import (
	"io"
	"strings"
	"testing"

	"github.com/alexfalkowski/tausch/internal/flag"
	"github.com/stretchr/testify/require"
)

func FuzzParseCommandName(f *testing.F) {
	f.Add("go", "version")
	f.Add(" test ", "my code")
	f.Add("-config", "cfg.yml")

	f.Fuzz(func(t *testing.T, first, second string) {
		if len(first)+len(second) > 4096 {
			t.Skip()
		}

		values, err := flag.Parse(io.Discard, []string{"--", first, second})
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(strings.Join([]string{first, second}, " ")), values.Name())
	})
}
