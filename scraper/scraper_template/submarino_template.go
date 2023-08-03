package scrapertemplate

import (
	"math/rand"
	"time"

	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/consts"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/utils"
	"go.uber.org/zap"
)

type defaultScraper[T model.ScrapingItems] struct {
	logger      *zap.SugaredLogger
	itemScraper ItemScraper[T]
}

func NewDefaultScraper[T model.ScrapingItems](logger *zap.Logger, itemScraper ItemScraper[T]) *defaultScraper[T] {
	s := &defaultScraper[T]{
		logger:      logger.Sugar(),
		itemScraper: itemScraper,
	}
	return s
}

func (ds *defaultScraper[T]) CollectData(baseURL string) ([]T, error) {
	limit := consts.DefaultLimit
	offset := uint(0)
	items := []T{}
	hasMoreItems := true
	currentItem := 0

	for hasMoreItems {
		listURL := utils.MountURL(baseURL, limit, offset)
		urls, totalItems, err := ds.itemScraper.CollectDetailURLs(listURL)
		if err != nil {
			return nil, err
		}

		ds.logger.Debugf("limit: %d | offset: %d | total: %d\n", limit, offset, totalItems)

		for _, url := range urls {
			currentItem += 1
			ds.logger.Debugf("current progress - item %d of %d\n", currentItem, totalItems)
			item, err := ds.itemScraper.CollectDetail(url)
			if err == nil {
				items = append(items, *item)
			}
			ds.randomDelay()
		}

		hasMoreItems = totalItems > (offset + limit)
		offset += limit
	}

	return items, nil
}

func (ds *defaultScraper[T]) randomDelay() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1 + rand.Intn(5) // n will be between 1 and 5
	for i := n; i > 0; i-- {
		ds.logger.Debugf("sleeping %d seconds ...\n", i)
		time.Sleep(time.Duration(1) * time.Second)
	}
}
