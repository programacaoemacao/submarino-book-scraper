package scraper

import (
	"github.com/programacaoemacao/submarino-book-scraper/model"
)

type AmazonBookScrapper interface {
	CollectData() ([]model.Book, error)
}
