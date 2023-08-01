package utils

import (
	"testing"

	"github.com/programacaoemacao/submarino-book-scraper/scraper/consts"
	"github.com/stretchr/testify/require"
)

func TestMountURL(t *testing.T) {
	t.Run("Mount a book category collect first page url", func(t *testing.T) {
		expectedURL := "https://www.submarino.com.br/categoria/livros/didaticos?limit=24&offset=0"
		url := MountURL("https://www.submarino.com.br/categoria/livros/didaticos", consts.DefaultLimit, 0)

		require.Equal(t, expectedURL, url)
	})
}
