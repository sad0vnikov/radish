package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/sad0vnikov/radish/redis"
)

//JSONFileConfigLoader is a struct for ConfigLoader implementation for loading config from JSON files
type JSONFileConfigLoader struct {
	Path string
}

//JSONContents represents JSON file format
type JSONContents struct {
	Servers   []redis.Server
	URLPrefix string
}

//Load config data from JSON file
func (jsonFile JSONFileConfigLoader) Load() (Config, error) {

	file, err := os.Open(jsonFile.Path)
	if err != nil {
		return Config{}, err
	}

	jsonContents := JSONContents{}

	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&jsonContents); err != nil {
		return config, err
	}

	c, err := fillConfig(jsonContents)
	config = c

	return config, err
}

func fillConfig(contents JSONContents) (Config, error) {
	config := Config{}
	config.Servers = make(map[string]redis.Server)
	for _, server := range contents.Servers {
		if _, prs := config.Servers[server.Name]; prs == true {
			return config, errors.New("server names should be unique in your config.json")
		}
		config.Servers[server.Name] = server
	}
	config.URLPrefix = contents.URLPrefix

	return config, nil
}
