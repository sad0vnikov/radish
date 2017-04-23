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

//HashKeyExists Hash has given key
func HashKeyExists(serverName, key, hashKey string) (bool, error) {

	conn, err := connector.GetByName(serverName)
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
func SetHashKey(serverName, key, hashKey, hashValue string) error {
	conn, err := connector.GetByName(serverName)
	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", key, hashKey, hashValue)
	if err != nil {
		return err
	}

	return nil
}

//DeleteHashValue deletes a Hash value
func DeleteHashValue(serverName, key, hashKey string) error {
	conn, err := connector.GetByName(serverName)
	if err != nil {
		return err
	}

	_, err = conn.Do("HDEL", key, hashKey)
	if err != nil {
		return err
	}

	return nil
}
