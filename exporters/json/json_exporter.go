package json

import (
	"encoding/json"
	"os"
)

type jsonExporter struct {
	file string
}

func NewJSONExporter(file string) *jsonExporter {
	return &jsonExporter{
		file: file,
	}
}

func (je *jsonExporter) Export(items interface{}) error {
	bytes, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("test.json", bytes, 0644)
	return err
}
