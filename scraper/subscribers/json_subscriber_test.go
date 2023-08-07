package subscribers

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func compareFiles(t *testing.T, filePath1 string, filePath2 string) {
	file1Bytes, err := os.ReadFile(filePath1)
	fmt.Printf("%s\n---------------------------------", string(file1Bytes))
	require.NoError(t, err)

	file2Bytes, err := os.ReadFile(filePath2)
	fmt.Printf("%s\n---------------------------------", string(file2Bytes))
	require.NoError(t, err)

	require.Equal(t, file1Bytes, file2Bytes)
}

func copyFile(t *testing.T, src string, dst string) {
	sourceFileStat, err := os.Stat(src)
	require.NoError(t, err)

	if !sourceFileStat.Mode().IsRegular() {
		t.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	require.NoError(t, err)

	defer source.Close()

	destination, err := os.Create(dst)
	require.NoError(t, err)

	defer destination.Close()
	_, err = io.Copy(destination, source)
	require.NoError(t, err)
}

func newLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestProcessData(t *testing.T) {
	t.Run("Test write a new json", func(t *testing.T) {
		dir := t.TempDir()
		jsonPath := path.Join(dir, "test.json")

		book := &model.Book{
			Title: "A random Book",
			Authors: []string{
				"Some author",
			},
			Rating: model.BookRating{
				Average:        5.0,
				TotalOfRatings: 12,
			},
		}

		subscriber := newJSONSubscriber[model.Book](jsonPath, newLogger())
		err := subscriber.ProcessData(book)

		require.NoError(t, err)
		require.FileExists(t, jsonPath)
		compareFiles(t, jsonPath, "./test_1.json")
	})

	t.Run("Test with an existing json", func(t *testing.T) {
		dir := t.TempDir()
		jsonPath := path.Join(dir, "test.json")
		copyFile(t, "./test_1.json", jsonPath)

		book := &model.Book{
			Title: "A cosmic Book",
			Authors: []string{
				"Another Author",
			},
			Rating: model.BookRating{
				Average:        4.1,
				TotalOfRatings: 10,
			},
		}

		subscriber := newJSONSubscriber[model.Book](jsonPath, newLogger())
		err := subscriber.ProcessData(book)

		require.NoError(t, err)
		require.FileExists(t, jsonPath)
		compareFiles(t, jsonPath, "./test_2.json")
	})
}
