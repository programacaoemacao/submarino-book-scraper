package main

import (
	exporter "github.com/programacaoemacao/submarino-book-scraper/exporters/interfaces"
	jsonexporter "github.com/programacaoemacao/submarino-book-scraper/exporters/json"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	bookscraper "github.com/programacaoemacao/submarino-book-scraper/scraper/items/book"
	scrapertemplate "github.com/programacaoemacao/submarino-book-scraper/scraper/scraper_template"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	defer logger.Sync() // flushes buffer, if any

	bookScraper := bookscraper.NewBookScraper(logger)
	scraperTemplate := scrapertemplate.NewDefaultScraper[model.Book](logger, bookScraper)

	books, err := scraperTemplate.CollectData("https://www.submarino.com.br/landingpage/trd-autoajuda?chave=trd-hi-at-generos-livros-blackfriday-autoajuda")
	if err != nil {
		logger.Sugar().Fatalf("can't collect all books data: %s", err.Error())
	}

	var exporter exporter.Exporter = jsonexporter.NewJSONExporter("./books.json")
	err = exporter.Export(books)
	if err != nil {
		logger.Sugar().Fatalf("can't collect all books data: %s", err.Error())
	}
}