package db

import (
	"errors"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
)

//Connections is an interface for objects storing Redis connections
type Connections interface {
	GetByName(serverName string, dbNum uint8) (redis.Conn, error)
	GetMaxDbNumsForServer(serverName string) (uint8, error)
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

//GetMaxDbNumsForServer returns a maxium db number for given Redis server
func (connections RedisConnections) GetMaxDbNumsForServer(serverName string) (uint8, error) {
	conn, err := connections.GetByName(serverName, 0)
	if err != nil {
		return 0, err
	}

	r, err := conn.Do("CONFIG", "GET", "DATABASES")
	cfgValues, err := redis.Strings(r, err)
	if err != nil {
		return 0, err
	}

	cnt, err := strconv.ParseUint(cfgValues[1], 10, 8)

	if err != nil {
		return 0, err
	}

	return uint8(cnt), nil
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

//GetMaxDbNumsForServer returns max db nums for a mocked connection
func (connections MockedConnections) GetMaxDbNumsForServer(serverName string) (uint8, error) {
	if connections.ConnectionMock == nil {
		return 0, errors.New("mock is not set, put your connection mock at ConectionMock struct field")
	}
	return uint8(8), nil
}

//GetServersWithConnectionData returns servers list from config with connection data like databases count, connection status, etc.
func GetServersWithConnectionData() map[string]rd.Server {
	configServers := config.Get().Servers
	result := make(map[string]rd.Server)
	for _, srv := range configServers {
		dbsCount, _ := GetMaxDbNumsForServer(srv.Name)
		srv.DatabasesCount = dbsCount
		result[srv.Name] = srv
	}

	return result
}
