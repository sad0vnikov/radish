package config

import "github.com/sad0vnikov/radish/redis"

//StubConfigLoader is a stubbed config loader for testing purposes
type StubConfigLoader struct {
}

//Load func loads stubbed app config
func (StubConfigLoader) Load() (Config, error) {
	servers := map[string]redis.Server{}
	servers["server1"] = redis.NewServer("server1", "127.0.0.1", 6379)
	servers["server2"] = redis.NewServer("server2", "127.0.0.1", 6380)
	servers["server3"] = redis.NewServer("server3", "127.0.0.1", 6381)

	config = Config{Servers: servers}

	return config, nil
}
