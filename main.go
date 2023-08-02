package main

import (
	"log"

	exporter "github.com/programacaoemacao/submarino-book-scraper/exporters/interfaces"
	jsonexporter "github.com/programacaoemacao/submarino-book-scraper/exporters/json"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/book"
	scraper "github.com/programacaoemacao/submarino-book-scraper/scraper/interfaces"
)

func main() {
	var bookScraper scraper.SubmarinoItemScraper[model.Book] = book.NewBookScraper()
	books, err := bookScraper.CollectData("https://www.submarino.com.br/landingpage/trd-livros-mais-vendidos")
	if err != nil {
		log.Fatalf("can't collect all books data: %s", err.Error())
	}

	var exporter exporter.Exporter = jsonexporter.NewJSONExporter("./books.json")
	err = exporter.Export(books)
	if err != nil {
		log.Fatalf("can't collect all books data: %s", err.Error())
	}
}
