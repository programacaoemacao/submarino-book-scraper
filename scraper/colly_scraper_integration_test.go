//go:build integration

package scraper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectData(t *testing.T) {
	t.Run("Collect URL's without error", func(t *testing.T) {
		scraper := NewBookScraper()

		// This url have only 2 pages
		baseURL := "https://www.submarino.com.br/categoria/livros/ciencias-exatas/f/loja-3p%7COlist+Store"
		books, err := scraper.CollectData(baseURL)

		require.NoError(t, err)
		require.NotNil(t, books)
	})
}
