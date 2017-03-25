package config

import "github.com/sad0vnikov/radish/redis"

//StubConfigLoader is a stubbed config loader for testing purposes
type StubConfigLoader struct {
}

func (StubConfigLoader) Load() (Config, error) {
	servers := make([]redis.Server, 3)
	servers[0] = redis.NewServer("server1", "127.0.0.1", 6379)
	servers[1] = redis.NewServer("server2", "127.0.0.1", 6380)
	servers[2] = redis.NewServer("server3", "127.0.0.1", 6381)

	config = Config{Servers: servers}

	return config, nil
}
