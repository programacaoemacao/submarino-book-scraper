package scraper

import "github.com/programacaoemacao/submarino-book-scraper/model"

type ScraperStrategy[T model.ScrapingItems] interface {
	CollectDetailURLs(itemsGridURL string) (urls []string, totalItems uint, err error)
	CollectDetail(detailURL string) (*T, error)
}

type ScraperSubscriber[T model.ScrapingItems] interface {
	ProcessData(item *T) error
}
