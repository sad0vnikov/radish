package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

type HashKey struct {
	key        string
	serverName string
}

//KeyType returns Hash key type
func (key HashKey) KeyType() string {
	return RedisHash
}

//PagesCount returns Hash key values pages count
func (key HashKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		return 0, err
	}

	r, err := conn.Do("HLEN", key.key)
	count, err := redis.Int(r, err)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return getValuesPagesCount(count, pageSize), err
}

//Values returns Hash key Values page
func (key HashKey) Values(pageNum int, pageSize int) (interface{}, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		return nil, err
	}

	var (
		cursor int64
		values []string
	)

	r, err := redis.Values(conn.Do("HSCAN", key.key, pageNum*pageSize, "COUNT", pageSize))
	r, err = redis.Scan(r, &cursor, &values)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	valuesMap := make(map[string]string)
	for i := 1; i < len(values); i = i + 2 {
		hashKey := values[i-1]
		hashValue := values[i]
		valuesMap[hashKey] = hashValue
	}

	return valuesMap, nil
}
