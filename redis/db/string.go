package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//StringKey is a key with redis String value
type StringKey struct {
	serverName string
	dbNum      uint8
	key        string
}

//Values returns StringKey values
func (key StringKey) Values(pageNum int, pageSize int) (interface{}, error) {
	return getStringKeyValue(key.serverName, key.dbNum, key.key)
}

//PagesCount returns StringKey pages count
func (StringKey) PagesCount(pageSize int) (int, error) {
	return 1, nil
}

//KeyType returns redis key type
func (StringKey) KeyType() string {
	return RedisString
}

//Set sets string value
func Set(serverName string, dbNum uint8, key, value string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		logger.Critical(err)
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//GetStringKeyValue returns a value for STRING type object
func getStringKeyValue(serverName string, dbNum uint8, key string) (RedisValue, error) {
	conn, err := connector.GetByName(serverName, dbNum)

	if err != nil {
		return RedisValue{}, fmt.Errorf("can't connect to server %v", serverName)
	}

	result, err := conn.Do("GET", key)
	if err != nil {
		logger.Error(err)
		return RedisValue{}, fmt.Errorf("can't get key %v value", key)
	}

	str, err := redis.String(result, err)
	value := RedisValue{
		Value:    str,
		IsBinary: isBinary(str),
	}
	return value, err
}
