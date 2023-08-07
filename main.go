package main

import (
	"os"

	"github.com/programacaoemacao/submarino-book-scraper/flags"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	scraper "github.com/programacaoemacao/submarino-book-scraper/scraper/scraper"
	bookstrategy "github.com/programacaoemacao/submarino-book-scraper/scraper/strategies/book"
	scrapersubscriber "github.com/programacaoemacao/submarino-book-scraper/scraper/subscribers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	defer logger.Sync()

	opts, err := flags.GetOptions(os.Args...)
	if err != nil {
		logger.Sugar().Fatal(err.Error())
	}

	bookStrategyScraper := bookstrategy.NewBookScraper(logger)
	submarinoScraper := scraper.NewDefaultScraper[model.Book](logger, bookStrategyScraper)

	outputSubscriber, err := scrapersubscriber.NewSubscriber[model.Book](opts.Output, logger)
	if err != nil {
		logger.Sugar().Fatalf("can't create optput subscriber: %s", err.Error())
	}

	subscribers := []scraper.ScraperSubscriber[model.Book]{
		outputSubscriber,
	}

	err = submarinoScraper.CollectData(opts.URLToCollect, subscribers)
	if err != nil {
		logger.Sugar().Fatalf("can't collect all books data: %s", err.Error())
	}
}
