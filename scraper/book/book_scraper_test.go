package book

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gocolly/colly"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/stretchr/testify/require"
)

func createNewMockBookScraper(t *testing.T) *bookScraper {
	scraper := NewBookScraper()
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(transport)

	scraper.collector = c
	return scraper
}

func getAbsoluteProjectRootDir(t *testing.T) string {
	dir, err := filepath.Abs("../..")
	require.NoError(t, err)
	return dir
}

func TestScrapeBooksURLS(t *testing.T) {
	t.Run("URLS have been collected", func(t *testing.T) {
		scraper := createNewMockBookScraper(t)
		projectRootDir := getAbsoluteProjectRootDir(t)

		url := "file://" + projectRootDir + "/test_files/example_books_page.html"

		// Fix this later
		expectedURLs := []string{
			"https:///produto/5092532919?pfm_index=1&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/128275610?pfm_index=2&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/5565139?pfm_index=3&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/5397698670?pfm_index=4&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/111489056?pfm_index=5&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/6025457120?pfm_index=6&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/2064089075?pfm_index=7&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/153630?pfm_index=8&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/121595813?pfm_index=9&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/132600243?pfm_index=10&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/4514117521?pfm_index=11&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/134289911?pfm_index=12&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/130207667?pfm_index=13&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/134495820?pfm_index=14&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/3292511021?pfm_index=15&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/1230296492?pfm_index=16&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/134494966?pfm_index=17&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/128275871?pfm_index=18&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/5144598478?pfm_index=19&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/134496662?pfm_index=20&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/124113761?pfm_index=21&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/9779533?pfm_index=22&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/4463939532?pfm_index=23&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
			"https:///produto/4801212100?pfm_index=24&pfm_page=category&pfm_pos=grid&pfm_type=category_page",
		}

		gottenURLs, totalItems, err := scraper.scrapeBooksURLS(url)
		require.NoError(t, err)
		require.NotZero(t, totalItems)
		require.Equal(t, gottenURLs, expectedURLs)
	})

	t.Run("Error when scraping - URL not found", func(t *testing.T) {
		scraper := createNewMockBookScraper(t)
		projectRootDir := getAbsoluteProjectRootDir(t)

		url := "file://" + projectRootDir + "/test_files/non_existent_file.html"
		urls, totalItems, err := scraper.scrapeBooksURLS(url)
		require.Nil(t, urls)
		require.Zero(t, totalItems)
		require.Error(t, err)
	})

	t.Run("Scraping - 0 total of items", func(t *testing.T) {
		scraper := createNewMockBookScraper(t)
		projectRootDir := getAbsoluteProjectRootDir(t)

		url := "file://" + projectRootDir + "/test_files/no_total_items.html"
		urls, totalItems, err := scraper.scrapeBooksURLS(url)
		require.Nil(t, urls)
		require.Zero(t, totalItems)
		require.Error(t, err)
	})
}

func TestScrapeBook(t *testing.T) {

	t.Run("Book have been collected", func(t *testing.T) {
		scraper := createNewMockBookScraper(t)
		projectRootDir := getAbsoluteProjectRootDir(t)

		expectedBook := &model.Book{
			CoverImageURL: "https://images-americanas.b2w.io/produtos/132332550/imagens/livro-as-4-disciplinas-da-execucao-garanta-o-foco-nas-metas-crucialmente-importantes/132332550_1_large.jpg",
			Title:         "as 4 disciplinas da execução: garanta o foco nas metas crucialmente importantes",
			Authors: []string{
				"Bill Moraes",
				"Chris McChesney",
				"Sean Covey",
			},
			Rating: model.BookRating{
				Average:        4.3,
				TotalOfRatings: 4,
			},
			PriceInCents:     9700,
			PaymentCondition: "em até 4x sem juros no cartão de crédito", // Get more info AND remove &nbsp; char
			Description:      "Você se recorda da última grande iniciativa que viu morrer na sua empresa?\nHouve algum estrondo? Ou foi sendo lenta e submissamente sufocada por outras prioridades? Quando desapareceu, provavelmente ninguém notou.\nO que aconteceu? O “redemoinho” de atividades urgentes necessário para manter as coisas funcionando no dia a dia devorou todo o tempo e energia que você precisava investir na execução da sua estratégia para o amanhã!\n\nAs 4 Disciplinas da Execução\u00a0(4DX) constituem uma fórmula simples, repetível e comprovada para você executar suas mais importantes prioridades estratégicas em meio ao redemoinho. Com Foco no crucialmente importante; Atuação nas medidas de direção; Manutenção de um placar envolvente e a Criação de uma cadência de responsabilidade, os líderes podem gerar resultados surpreendentes, até mesmo quando a execução da estratégia demanda significativa mudança no comportamento de suas equipes.\n4DX não é teoria. É um conjunto de práticas comprovadas, que já foram testadas e aperfeiçoadas por centenas de organizações e milhares de equipes nos últimos 10 anos. Quando empresas ou indivíduos aderem a estas disciplinas, alcançam ótimos resultados, independentemente do objetivo a ser alcançado. 4DX representa um novo modo de pensar e trabalhar essencial para a prosperidade no clima competitivo da atualidade. \nTrata-se de um livro que nenhum líder pode deixar de ler.",
			Metadata: model.BookMetadata{
				Pages:       352,
				Languages:   []string{"Português", "Inglês"},
				Publisher:   "Alta Books",
				PublishDate: "12/06/2017",
				ISBN10:      "8550801399",
				ISBN13:      "9788550801391",
				Edition:     "1° Ed.",
			},
		}

		url := "file://" + projectRootDir + "/test_files/example_book_1.html"
		gottenBook, err := scraper.scrapeBook(url)
		require.NoError(t, err)
		require.Equal(t, expectedBook, gottenBook)
	})

	t.Run("Error when scraping - URL not found", func(t *testing.T) {
		scraper := createNewMockBookScraper(t)
		projectRootDir := getAbsoluteProjectRootDir(t)

		url := "file://" + projectRootDir + "/test_files/non_existent_file.html"
		book, err := scraper.scrapeBook(url)
		require.Nil(t, book)
		require.Error(t, err)
	})
}
