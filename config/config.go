package config

import (
	"github.com/sad0vnikov/radish/redis"
)

//Config is struct which stores configuration data
type Config struct {
	Servers map[string]redis.Server
}

//Loader is an interface for configuration loading logic
type Loader interface {
	Load() (Config, error)
}

var config = Config{}

//Get config data
func Get() Config {
	return config
}
