package subscribers

import (
	"errors"
	"strings"

	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/scraper"
	"go.uber.org/zap"
)

func NewSubscriber[T model.ScrapingItems](filePath string, logger *zap.Logger) (scraper.ScraperSubscriber[T], error) {
	switch {
	case strings.HasSuffix(filePath, ".json"):
		return newJSONSubscriber[T](filePath, logger), nil
	default:
		return nil, errors.New("file extension not supported")
	}
}
