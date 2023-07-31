package scraper

import "github.com/programacaoemacao/submarino-book-scraper/model"

type Item interface {
	model.Book
}

type SubmarinoItemScraper[T Item] interface {
	CollectData(baseURL string) ([]T, error)
}
