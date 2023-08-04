package main

import (
	"github.com/programacaoemacao/submarino-book-scraper/model"
	scraper "github.com/programacaoemacao/submarino-book-scraper/scraper/scraper"
	bookstrategy "github.com/programacaoemacao/submarino-book-scraper/scraper/strategies/book"
	subscriber "github.com/programacaoemacao/submarino-book-scraper/scraper/subscribers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	defer logger.Sync()

	bookStrategyScraper := bookstrategy.NewBookScraper(logger)
	submarinoScraper := scraper.NewDefaultScraper[model.Book](logger, bookStrategyScraper)

	url := "https://www.submarino.com.br/landingpage/trd-autoajuda?chave=trd-hi-at-generos-livros-blackfriday-autoajuda"
	subscribers := []scraper.ScraperSubscriber[model.Book]{
		subscriber.NewJSONSubscriber[model.Book]("books.json", logger),
	}

	err := submarinoScraper.CollectData(url, subscribers)
	if err != nil {
		logger.Sugar().Fatalf("can't collect all books data: %s", err.Error())
	}
}
