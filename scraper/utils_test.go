package scraper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMountURL(t *testing.T) {
	t.Run("Mount a book category collect first page url", func(t *testing.T) {
		expectedURL := "https://www.submarino.com.br/categoria/livros/didaticos?limit=24&offset=0"
		url := mountURL("https://www.submarino.com.br/categoria/livros/didaticos", defaultLimit, 0)

		require.Equal(t, expectedURL, url)
	})
}
