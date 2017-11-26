package config

import (
	"github.com/sad0vnikov/radish/redis"
)

//Config is struct which stores configuration data
type Config struct {
	HttpPort int16
	Servers   map[string]redis.Server
	URLPrefix string
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
