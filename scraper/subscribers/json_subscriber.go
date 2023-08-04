package subscribers

import (
	"encoding/json"
	"os"

	"github.com/programacaoemacao/submarino-book-scraper/model"
	"go.uber.org/zap"
)

type jsonSubscriber[T model.ScrapingItems] struct {
	filePath string
	logger   *zap.SugaredLogger
}

func NewJSONSubscriber[T model.ScrapingItems](filePath string, logger *zap.Logger) *jsonSubscriber[T] {
	return &jsonSubscriber[T]{
		filePath: filePath,
		logger:   logger.Sugar(),
	}
}

func (j *jsonSubscriber[T]) ProcessData(item *T) error {
	j.logger.Debugf("opening file %q", j.filePath)
	jsonFile, err := os.OpenFile(j.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		j.logger.Errorf("error at opening file %q: %s", j.filePath, err.Error())
		return err
	}

	defer func() {
		j.logger.Debugf("closing file %q", j.filePath)
		jsonFile.Close()
	}()

	fileInfo, err := jsonFile.Stat()
	if err != nil {
		j.logger.Errorf("error getting file %q info: %s", j.filePath, err.Error())
		return err
	}

	if fileInfo.Size() > 0 {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			j.logger.Errorf("error transforming item to json content: %s", err.Error())
			return err
		}

		itemBytes = append([]byte(","), itemBytes...)
		itemBytes = append(itemBytes, []byte("]")...)

		// Removing the last char
		err = jsonFile.Truncate(fileInfo.Size() - 1)
		if err != nil {
			j.logger.Errorf("error truncating last char of file %q", j.filePath)
			return err
		}

		j.logger.Debugf("writing content on file %q", j.filePath)
		_, err = jsonFile.Write(itemBytes)
		if err != nil {
			j.logger.Errorf("error appending content to file %q", j.filePath)
			return err
		}
	} else {
		sliceItems := []T{*item}
		itemBytes, err := json.Marshal(sliceItems)
		if err != nil {
			j.logger.Errorf("error transforming item to json content: %s", err.Error())
			return err
		}

		j.logger.Debugf("writing content on file %q", j.filePath)
		_, err = jsonFile.Write(itemBytes)
		if err != nil {
			j.logger.Errorf("error appending content to file %q", j.filePath)
			return err
		}
	}

	err = jsonFile.Sync()
	if err != nil {
		j.logger.Errorf("error syncing file %q", j.filePath)
		return err
	}

	return nil
}
