package scraper

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gocolly/colly"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/stretchr/testify/assert"
)

func TestScrapeBook(t *testing.T) {

	t.Run("Happy path -> Page was fetched", func(t *testing.T) {
		scrapper := NewScraper()
		transport := &http.Transport{}
		transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

		c := colly.NewCollector()
		c.WithTransport(transport)

		scrapper.collector = c

		os.Chdir("..")
		dir, err := filepath.Abs("")
		if err != nil {
			panic(err)
		}

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

		url := "file://" + dir + "/test_files/example_book_1.html"
		gottenBook, err := scrapper.scrapeBook(url)
		assert.NoError(t, err)
		assert.Equal(t, expectedBook, gottenBook)
	})
}
