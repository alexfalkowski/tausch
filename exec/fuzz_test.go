package exec_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/tausch/exec"
	"github.com/stretchr/testify/require"
)

func FuzzCommandPrefixesDelimiter(f *testing.F) {
	f.Add("go", "version", "")
	f.Add("go", "env", "GOPATH")
	f.Add("-version", "--json", "")

	f.Fuzz(func(t *testing.T, name, first, second string) {
		if len(name)+len(first)+len(second) > 4096 {
			t.Skip()
		}

		want := []string{"--", name, first, second}

		command := exec.Command(name, first, second)
		require.Equal(t, want, command.Args[1:])

		commandContext := exec.CommandContext(context.Background(), name, first, second)
		require.Equal(t, want, commandContext.Args[1:])
	})
}
