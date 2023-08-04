package main

import (
	"github.com/programacaoemacao/submarino-book-scraper/model"
	bookscraper "github.com/programacaoemacao/submarino-book-scraper/scraper/items/book"
	scrapertemplate "github.com/programacaoemacao/submarino-book-scraper/scraper/scraper_template"
	subscriber "github.com/programacaoemacao/submarino-book-scraper/scraper/subscribers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	defer logger.Sync()

	bookScraper := bookscraper.NewBookScraper(logger)
	scraperTemplate := scrapertemplate.NewDefaultScraper[model.Book](logger, bookScraper)

	url := "https://www.submarino.com.br/landingpage/trd-autoajuda?chave=trd-hi-at-generos-livros-blackfriday-autoajuda"
	subscribers := []scrapertemplate.ScraperSubscriber[model.Book]{
		subscriber.NewJSONSubscriber[model.Book]("books.json", logger),
	}
	err := scraperTemplate.CollectData(url, subscribers)
	if err != nil {
		logger.Sugar().Fatalf("can't collect all books data: %s", err.Error())
	}
}
