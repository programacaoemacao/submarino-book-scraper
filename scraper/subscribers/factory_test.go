package subscribers

import (
	"testing"

	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewSubscriber(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	t.Run("Wrong file format", func(t *testing.T) {
		subscriber, err := NewSubscriber[model.Book]("somepath.xml", logger)
		require.Nil(t, subscriber)
		require.ErrorContains(t, err, "file extension not supported")
	})

	t.Run("JSON Output", func(t *testing.T) {
		subscriber, err := NewSubscriber[model.Book]("somepath.json", logger)
		require.IsType(t, (*jsonSubscriber[model.Book])(nil), subscriber)
		require.NoError(t, err)
	})
}
