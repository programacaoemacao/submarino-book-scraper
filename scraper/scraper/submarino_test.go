package scraper

import (
	"testing"

	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockItemsScraper struct {
}

func (m *mockItemsScraper) CollectDetailURLs(itemsGridURL string) (urls []string, totalItems uint, err error) {
	urls = []string{
		"http://test.com",
	}
	return urls, 1, nil
}

func (m *mockItemsScraper) CollectDetail(detailURL string) (*model.Book, error) {
	book := &model.Book{
		CoverImageURL: "https://images-americanas.b2w.io/produtos/3097551545/imagens/livro-como-superar-seus-limites-internos-aprenda-a-vencer-seus-bloqueios-e-suas-batalhas-interiores-de-criatividade/3097551545_1_large.jpg",
		Title:         "como superar seus limites internos: aprenda a vencer seus bloqueios e suas batalhas interiores de criatividade",
		Authors: []string{
			"Steven Pressfield",
		},
		Rating: model.BookRating{
			Average:        4.1,
			TotalOfRatings: 46,
		},
		PriceInCents:     3168,
		PaymentCondition: "ver mais sugestões",
		Description:      "Em Como Superar seus Limites Internos - nova edição do clássico A Guerra da Arte -, o romancista best-seller Steven Pressfield identifica o inimigo que todos precisamos enfrentar em nós mesmos, traçando um plano de batalha para o vencermos e apresentando importantes ensinamentos para alcançarmos o máximo de sucesso.\nEle enfatiza ainda a resolução necessária para reconhecer e superar os obstáculos à ambição, e mostra, com clareza, como chegar ao mais alto nível de disciplina criativa. Com prefácio exclusivo de Lúcia Helena Galvão, professora de Filosofia da organização Nova Acrópole do Brasil há 31 anos, este livro é simplesmente A Arte da Guerra de Sun Tzu para a alma.",
		Metadata: model.BookMetadata{
			Pages: 200,
			Languages: []string{
				"Português",
			},
			Publisher:   "Cultrix",
			PublishDate: "05/05/2021",
			ISBN10:      "6557360973",
			ISBN13:      "9786557360972",
			Edition:     "1° Ed.",
		},
	}
	return book, nil
}

func newMockedScraperStrategy(t *testing.T) *defaultScraper[model.Book] {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	ScraperStrategy := &mockItemsScraper{}

	return NewDefaultScraper[model.Book](logger, ScraperStrategy)
}

func TestCollectData(t *testing.T) {
	// TODO: Inject the delay function to avoid tests to wait 4 or 5 seconds to run
	t.Run("Success", func(t *testing.T) {
		ScraperStrategy := newMockedScraperStrategy(t)

		subscribers := []ScraperSubscriber[model.Book]{}
		err := ScraperStrategy.CollectData("http://test.com", subscribers)
		require.NoError(t, err)
	})
}
