package data_sources

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONFileSource struct {
	FilePath string
}

// FetchData reads data from a JSON file
func (j *JSONFileSource) FetchData() ([]map[string]interface{}, error) {
	file, err := os.ReadFile(j.FilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("could not parse JSON file: %v", err)
	}

	return data, nil
}
