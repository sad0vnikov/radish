package db

import (
	"errors"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/logger"
)

//Connections is an interface for objects storing Redis connections
type Connections interface {
	GetByName(serverName string, dbNum uint8) (redis.Conn, error)
}

//RedisConnections is a struct storing redis connection
type RedisConnections struct {
	pool *redis.Pool
}

//GetByName function returns redigo Redis connection instance by server name
func (connections RedisConnections) GetByName(serverName string, dbNum uint8) (redis.Conn, error) {
	server, prs := config.Get().Servers[serverName]
	if !prs {
		return nil, errors.New("no server with name " + serverName + " found")
	}

	if connections.pool == nil {
		connections.pool = &redis.Pool{
			MaxIdle: 3,
			Dial: func() (redis.Conn, error) {
				serverAddr := server.Host + ":" + strconv.Itoa(server.Port)
				logger.Info("connecting to Redis server " + serverAddr)
				conn, err := redis.Dial("tcp", serverAddr)
				logger.Info("connected to Redis server " + serverAddr)
				return conn, err
			},
		}
	}

	var c redis.Conn
	c = connections.pool.Get()
	c.Do("SELECT", dbNum)

	return c, nil
}

//MockedConnections is a struct storing mocked redis connections
type MockedConnections struct {
	ConnectionMock redis.Conn
}

//GetByName returns mocked connection
func (connections MockedConnections) GetByName(serverName string) (redis.Conn, error) {
	if connections.ConnectionMock == nil {
		return nil, errors.New("mock is not set, put your connection mock at ConectionMock struct field")
	}
	return connections.ConnectionMock, nil
}
