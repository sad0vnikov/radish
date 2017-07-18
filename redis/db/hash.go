package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
)

type HashKey struct {
	key        string
	serverName string
	dbNum      uint8
}

//KeyType returns Hash key type
func (key HashKey) KeyType() string {
	return RedisHash
}

//PagesCount returns Hash key values pages count
func (key HashKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName, key.dbNum)
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
	conn, err := connector.GetByName(key.serverName, key.dbNum)
	if err != nil {
		return nil, err
	}

	var (
		values []string
	)

	r, err := conn.Do("HGETALL", key.key)
	values, err = redis.Strings(r, err)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	offsetStart, offsetEnd, err := rd.GetPageRangeForStrings(values, pageSize*2, pageNum)

	if err != nil {
		return nil, err
	}
	valuesPage := values[offsetStart:offsetEnd]

	valuesMap := make(map[string]RedisValue)
	for i := 1; i < len(valuesPage); i = i + 2 {
		hashKey := valuesPage[i-1]
		hashValue := valuesPage[i]
		valuesMap[hashKey] = RedisValue{
			Value:    hashValue,
			IsBinary: isBinary(hashValue),
		}
	}

	return valuesMap, nil
}

//HashKeyExists Hash has given key
func HashKeyExists(serverName string, dbNum uint8, key, hashKey string) (bool, error) {

	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return false, err
	}

	r, err := conn.Do("HEXISTS", key, hashKey)
	exists, err := redis.Bool(r, err)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	return exists, nil
}

//SetHashKey sets a hash value
func SetHashKey(serverName string, dbNum uint8, key, hashKey, hashValue string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", key, hashKey, hashValue)
	if err != nil {
		return err
	}

	return nil
}

//UpdateHashKey updates a value and hash key
func UpdateHashKey(serverName string, dbNum uint8, key, hashKey, newHashKey, hashValue string) error {
	if hashKey == newHashKey {
		return SetHashKey(serverName, dbNum, key, hashKey, hashValue)
	}

	ex, err := HashKeyExists(serverName, dbNum, key, newHashKey)
	if err != nil {
		return err
	}

	if ex {
		return fmt.Errorf("hash key %s already exists in hash %s", newHashKey, key)
	}

	conn, err := connector.GetByName(serverName, dbNum)
	_, err = conn.Do("MULTI")
	if err != nil {
		return err
	}

	DeleteHashValue(serverName, dbNum, key, hashKey)
	SetHashKey(serverName, dbNum, key, newHashKey, hashValue)

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

//DeleteHashValue deletes a Hash value
func DeleteHashValue(serverName string, dbNum uint8, key, hashKey string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	_, err = conn.Do("HDEL", key, hashKey)
	if err != nil {
		return err
	}

	return nil
}
