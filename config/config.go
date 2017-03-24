package config

import (
	"encoding/json"
	"os"

	"../redis"
)

type Config struct {
	Servers []redis.RedisServer
}

var config = Config{}

func Load() (Config, error) {

	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}

	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func Get() Config {
	return config
}
