package config

import (
	"encoding/json"
	"os"
)

type Transformation struct {
	Type      string            `json:"type"`
	Fields    map[string]string `json:"fields,omitempty"`
	Condition string            `json:"condition,omitempty"`
}

type Config struct {
	Transformations []Transformation `json:"transformations"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
