package unicode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveUnicodeChar(t *testing.T) {
	t.Run("Test with \u00A0 char", func(t *testing.T) {
		input := "Hello\u00A0I like pop-soda"
		expectedOutput := "Hello I like pop-soda"

		fixedText := replaceUnicodeChar(input, '\u00A0', ' ')
		require.Equal(t, expectedOutput, fixedText)
	})

	t.Run("Test with \u2013 char", func(t *testing.T) {
		input := "Naruto\u2013kun"
		expectedOutput := "Naruto-kun"

		fixedText := replaceUnicodeChar(input, '\u2013', '-')
		require.Equal(t, expectedOutput, fixedText)
	})
}

func TestNormalize(t *testing.T) {
	t.Run("Test with \u00A0 and \u2013 char", func(t *testing.T) {
		input := "Hello\u00A0Naruto\u2013kun"
		expectedOutput := "Hello Naruto-kun"

		fixedText := Normalize(input)
		require.Equal(t, expectedOutput, fixedText)
	})
}
