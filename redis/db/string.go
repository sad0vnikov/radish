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

//StringKeyValue stores String key vInfo data
type StringKeyValue struct {
	value       RedisValue
	valueLoaded bool
	query       *KeyValuesQuery
	key         StringKey
}

//Values returns StringKey vInfo
func (values *StringKeyValue) Values() (interface{}, error) {

	if !values.valueLoaded {
		value, err := getStringKeyValue(values.key.serverName, values.key.dbNum, values.key.key)
		if err != nil {
			return nil, err
		}
		values.valueLoaded = true
		values.value = value
	}

	return values.value, nil
}

//PagesCount returns StringKey vInfo
func (values *StringKeyValue) PagesCount() (int, error) {
	return 1, nil
}

//Values returns String vInfo object
func (key StringKey) Values(query *KeyValuesQuery) KeyValues {
	return &StringKeyValue{query: query, key: key}
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
		return RedisValue{}, fmt.Errorf("can't connect to server %vInfo", serverName)
	}

	result, err := conn.Do("GET", key)
	if err != nil {
		logger.Error(err)
		return RedisValue{}, fmt.Errorf("can't get key %vInfo value", key)
	}

	str, err := redis.String(result, err)
	value := RedisValue{
		Value:    str,
		IsBinary: isBinary(str),
	}
	return value, err
}
