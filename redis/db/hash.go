package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
)

type HashValues struct {
	values           map[string]RedisValue
	pagesCount       int
	valuesLoaded     bool
	pagesCountLoaded bool
	query            *KeyValuesQuery
	key              HashKey
}

func (values *HashValues) Values() (interface{}, error) {
	if !values.valuesLoaded {
		err := values.loadValues()
		if err != nil {
			return nil, err
		}
	}

	return values.values, nil
}

func (values *HashValues) PagesCount() (int, error) {
	if !values.pagesCountLoaded {
		err := values.calculatePagesCount()
		if err != nil {
			return 0, err
		}
	}

	return values.pagesCount, nil
}

type HashKey struct {
	key        string
	serverName string
	dbNum      uint8
}

//KeyType returns Hash key type
func (key HashKey) KeyType() string {
	return RedisHash
}

func (key HashKey) Values(query *KeyValuesQuery) KeyValues {
	return &HashValues{
		query: query,
		key:   key,
	}
}

//PagesCount returns Hash key vInfo pages count
func (vInfo *HashValues) calculatePagesCount() error {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return err
	}

	count := 0
	if vInfo.query.Mask == "*" {
		r, err := conn.Do("HLEN", vInfo.key.key)
		count, err = redis.Int(r, err)
	} else {
		if !vInfo.valuesLoaded {
			err = vInfo.loadValues()
		}
		count = len(vInfo.values)
	}

	if err != nil {
		logger.Error(err)
		return err
	}

	vInfo.pagesCount = count
	return nil
}

//Values returns Hash key Values page
func (vInfo *HashValues) loadValues() error {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return err
	}

	var (
		values []string
	)

	r, err := conn.Do("HGETALL", vInfo.key.key)
	values, err = redis.Strings(r, err)
	if err != nil {
		logger.Error(err)
		return err
	}

	offsetStart, offsetEnd, err := rd.GetPageRangeForStrings(values, vInfo.query.PageSize*2, vInfo.query.PageNum)

	if err != nil {
		return err
	}
	valuesPage := values[offsetStart:offsetEnd]

	valuesMap := make(map[string]RedisValue)
	for i := 1; i < len(valuesPage); i = i + 2 {
		hashKey := valuesPage[i-1]
		hashValue := valuesPage[i]
		if matchStringValueWithMask(hashKey, vInfo.query.Mask) {
			valuesMap[hashKey] = RedisValue{
				Value:    hashValue,
				IsBinary: isBinary(hashValue),
			}
		}
	}

	vInfo.values = valuesMap
	return nil
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
