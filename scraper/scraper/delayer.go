package scraper

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type delayer interface {
	delay()
}

type randomDelayer struct {
	logger *zap.SugaredLogger
}

func newRandomDelayer(logger *zap.Logger) *randomDelayer {
	return &randomDelayer{
		logger: logger.Sugar(),
	}
}

func (rd *randomDelayer) delay() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1 + rand.Intn(5) // n will be between 1 and 5
	for i := n; i > 0; i-- {
		rd.logger.Debugf("sleeping %d seconds ...\n", i)
		time.Sleep(time.Duration(1) * time.Second)
	}
}
