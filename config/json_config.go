package config

import (
	"encoding/json"
	"os"
)

//JSONFileConfigLoader is a struct for ConfigLoader implementation for loading config from JSON files
type JSONFileConfigLoader struct {
	Path string
}

//Load config data from JSON file
func (jsonFile JSONFileConfigLoader) Load() (Config, error) {

	file, err := os.Open(jsonFile.Path)
	if err != nil {
		return Config{}, err
	}

	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
