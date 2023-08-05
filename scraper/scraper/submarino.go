package scraper

import (
	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/consts"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/utils"
	"go.uber.org/zap"
)

type defaultScraper[T model.ScrapingItems] struct {
	logger          *zap.SugaredLogger
	scraperStrategy ScraperStrategy[T]
	delayer         delayer
}

func NewDefaultScraper[T model.ScrapingItems](logger *zap.Logger, scraperStrategy ScraperStrategy[T]) *defaultScraper[T] {
	s := &defaultScraper[T]{
		logger:          logger.Sugar(),
		scraperStrategy: scraperStrategy,
		delayer:         newRandomDelayer(logger),
	}
	return s
}

func (ds *defaultScraper[T]) CollectData(baseURL string, subscribers []ScraperSubscriber[T]) error {
	limit := consts.DefaultLimit
	offset := uint(0)
	hasMoreItems := true
	currentItem := 0

	for hasMoreItems {
		listURL := utils.MountURL(baseURL, limit, offset)
		urls, totalItems, err := ds.scraperStrategy.CollectDetailURLs(listURL)
		// TODO: Implement a better error treatment
		if err != nil {
			ds.logger.Errorf("error getting urls to collect: %s", err.Error())
			return err
		}

		ds.logger.Debugf("limit: %d | offset: %d | total: %d", limit, offset, totalItems)

		for _, url := range urls {
			currentItem += 1
			ds.logger.Debugf("current progress - item %d of %d", currentItem, totalItems)
			item, err := ds.scraperStrategy.CollectDetail(url)
			if err == nil {
				// Error supressed for simplicity
				_ = ds.notifySubscribers(subscribers, item)
			}
			ds.delayer.delay()
		}

		hasMoreItems = totalItems > (offset + limit)
		offset += limit
	}

	return nil
}

func (ds *defaultScraper[T]) notifySubscribers(subscribers []ScraperSubscriber[T], item *T) error {
	for _, s := range subscribers {
		// Error supressed for simplicity
		_ = s.ProcessData(item)
	}
	return nil
}
