package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
)

//SetKey is a redis SET key
type SetKey struct {
	serverName string
	dbNum      uint8
	key        string
}

//KeyType returns Redis SET type
func (key SetKey) KeyType() string {
	return RedisSet
}

//PagesCount returns Redis SET values pages count
func (key SetKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName, key.dbNum)
	if err != nil {
		return 0, err
	}

	r, err := conn.Do("SCARD", key.key)
	count, err := redis.Int(r, err)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return getValuesPagesCount(count, pageSize), nil
}

//Values returns a SET values page
func (key SetKey) Values(pageNum int, pageSize int) (interface{}, error) {
	conn, err := connector.GetByName(key.serverName, key.dbNum)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var (
		values []string
	)

	r, err := conn.Do("SMEMBERS", key.key)
	values, err = redis.Strings(r, err)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	offsetStart, offsetEnd, err := rd.GetPageRangeForStrings(values, pageSize, pageNum)

	if err != nil {
		return nil, err
	}
	valuesPageStrings := values[offsetStart:offsetEnd]

	valuesPage := make([]RedisValue, len(valuesPageStrings))
	for i, s := range valuesPageStrings {
		valuesPage[i] = RedisValue{
			Value:    s,
			IsBinary: isBinary(s),
		}
	}

	return valuesPage, nil
}

//AddValueToSet adds a new member to a set
func AddValueToSet(serverName string, dbNum uint8, key, value string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	_, err = conn.Do("SADD", key, value)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//UpdateSetValue updates a set member
func UpdateSetValue(serverName string, dbNum uint8, key, value, newValue string) error {
	err := DeleteSetValue(serverName, dbNum, key, value)
	if err != nil {
		return err
	}

	err = AddValueToSet(serverName, dbNum, key, newValue)
	if err != nil {
		return err
	}
	return nil
}

//DeleteSetValue removes a set member
func DeleteSetValue(serverName string, dbNum uint8, key, value string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	_, err = conn.Do("SREM", key, value)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
