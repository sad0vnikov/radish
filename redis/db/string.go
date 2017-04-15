package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//StringKey is a key with redis String value
type StringKey struct {
	serverName string
	key        string
}

//Values returns StringKey values
func (key StringKey) Values(pageNum, pageSize int) (interface{}, error) {
	return getStringKeyValue(key.serverName, key.key)
}

//PagesCount returns StringKey pages count
func (StringKey) PagesCount(pageSize int) (int, error) {
	return 1, nil
}

//KeyType returns redis key type
func (StringKey) KeyType() string {
	return RedisString
}

//GetStringKeyValue returns a value for STRING type object
func getStringKeyValue(serverName, key string) (string, error) {
	conn, err := connector.GetByName(serverName)

	if err != nil {
		return "", fmt.Errorf("can't connect to server %v", serverName)
	}

	result, err := conn.Do("GET", key)
	if err != nil {
		logger.Error(err)
		return "", fmt.Errorf("can't get key %v value", key)
	}

	value, err := redis.String(result, err)
	return value, err
}
