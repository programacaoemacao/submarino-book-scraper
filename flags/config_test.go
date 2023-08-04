package flags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetOptions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Cleaning other flags
		os.Args = []string{os.Args[0]}
		os.Args = append(os.Args, "--url=http://test.com")
		os.Args = append(os.Args, "--output=teste.json")

		opts, err := GetOptions(os.Args...)
		require.NoError(t, err)
		require.Equal(t, "http://test.com", opts.URLToCollect)
		require.Equal(t, "teste.json", opts.Output)
	})

	t.Run("Error", func(t *testing.T) {
		opts, err := GetOptions([]string{"-t=test"}...)
		require.ErrorContains(t, err, "unknown flag `t'")
		require.Nil(t, opts)
	})
}
