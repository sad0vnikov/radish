package db

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
	"regexp"
	"strconv"
	"strings"
)

//Connections is an interface for objects storing Redis connections
type Connections interface {
	GetByName(serverName string, dbNum uint8) (redis.Conn, error)
	GetMaxDbNumsForServer(serverName string) (uint8, error)
	GetServerKeyspaceStat(serverName string) (map[string]rd.ServerKeyspaceStat, error)
	GetServerStat(serverName string) (rd.ServerStat, error)
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

func (connections RedisConnections) GetServerStat(serverName string) (rd.ServerStat, error) {
	conn, err := connections.GetByName(serverName, 0)
	if err != nil {
		return rd.ServerStat{}, err
	}

	r, err := conn.Do("INFO")
	info, err := redis.String(r, err)
	if err != nil {
		return rd.ServerStat{}, err
	}

	return parseInfoStrings(strings.Split(info, "\r\n")), nil
}

func (connections RedisConnections) GetServerKeyspaceStat(serverName string) (map[string]rd.ServerKeyspaceStat, error) {
	conn, err := connections.GetByName(serverName, 0)
	if err != nil {
		return nil, err
	}

	r, err := conn.Do("INFO", "keyspace")
	rs, err := redis.String(r, err)
	info := strings.Split(rs, "\r\n")

	databaseStatRowsCount := len(info) - 1
	stat := make(map[string]rd.ServerKeyspaceStat, databaseStatRowsCount)
	dbNum := 0
	for i := 1; i < len(info); i++ {
		dbName := "db" + strconv.Itoa(dbNum)
		statString := info[i]

		dbKeyspaceStat := parseKeyspaceStatString(statString)
		dbNum++

		stat[dbName] = dbKeyspaceStat
	}
	return stat, nil
}

func parseKeyspaceStatString(statString string) rd.ServerKeyspaceStat {
	parseKeyspaceInfoRegexp := "[,:]?([a-z_]+=[a-z0-9]+)"
	r := regexp.MustCompile(parseKeyspaceInfoRegexp)
	stats := r.FindAllString(statString, -1)
	dbKeyspaceStat := rd.ServerKeyspaceStat{}
	for _, statString := range stats {
		splittedStat := strings.Split(statString, "=")
		statName := splittedStat[0][1:]
		statValue := splittedStat[1]
		switch statName {
		case "keys":
			intValue, _ := strconv.ParseInt(statValue, 10, 64)
			dbKeyspaceStat.KeysCount = intValue
		}
	}
	return dbKeyspaceStat
}

func parseInfoStrings(info []string) rd.ServerStat {
	c := rd.ServerStat{}

	for _, s := range info {
		if !strings.Contains(s, ":") {
			continue
		}

		splittedRow := strings.Split(s, ":")
		statName := splittedRow[0]
		statValue := splittedRow[1]
		intValue, _ := strconv.ParseInt(statValue, 10, 64)
		switch statName {
		case "used_memory_human":
			c.UsedMemoryHuman = statValue
		case "used_memory":
			c.UsedMemoryBytes = intValue
		case "max_memory_human":
			c.MaxMemoryHuman = statValue
		case "max_memory":
			c.MaxMemoryBytes = intValue
		case "redis_version":
			c.RedisVersion = statValue
		case "uptime_in_seconds":
			c.UptimeInSeconds = intValue
		case "connected_clients":
			c.ConnectedClientsCount = intValue
		}
	}
	return c
}

//MockedConnections is a struct storing mocked redis connections
type MockedConnections struct {
	ConnectionMock redis.Conn
}

//GetByName returns mocked connection
func (connections MockedConnections) GetByName(serverName string, dbNum uint8) (redis.Conn, error) {
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

func (connections MockedConnections) GetServerStat(serverName string) (rd.ServerStat, error) {
	return rd.ServerStat{}, nil
}

func (connections MockedConnections) GetServerKeyspaceStat(serverName string) (map[string]rd.ServerKeyspaceStat, error) {
	return map[string]rd.ServerKeyspaceStat{}, nil
}

//GetServersWithConnectionData returns servers list from config with connection data like databases count, connection status, etc.
func GetServersWithConnectionData() map[string]rd.Server {
	configServers := config.Get().Servers
	result := make(map[string]rd.Server)
	for _, srv := range configServers {
		dbsCount, err := GetMaxDbNumsForServer(srv.Name)
		if err == nil {
			srv.ConnectionCheckPassed = true
		}
		srv.DatabasesCount = dbsCount
		if srv.ConnectionCheckPassed {
			serverStat, err := connector.GetServerStat(srv.Name)
			if err != nil {
				logger.Error(err)
			}
			keyspaceStat, err := connector.GetServerKeyspaceStat(srv.Name)
			if err != nil {
				logger.Error(err)
			}
			srv.ServerStat = serverStat
			srv.KeyspaceStat = keyspaceStat
		}

		result[srv.Name] = srv
	}

	return result
}
