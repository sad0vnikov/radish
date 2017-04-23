package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//SetKey is a redis SET key
type SetKey struct {
	serverName string
	key        string
}

//KeyType returns Redis SET type
func (key SetKey) KeyType() string {
	return RedisSet
}

//PagesCount returns Redis SET values pages count
func (key SetKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName)
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
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var (
		cursor int64
		values []string
	)

	r, err := redis.Values(conn.Do("SSCAN", key.key, pageNum*pageSize, "COUNT", pageSize))
	r, err = redis.Scan(r, &cursor, &values)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return values, nil
}

//AddValueToSet adds a new member to a set
func AddValueToSet(serverName, key, value string) error {
	conn, err := connector.GetByName(serverName)
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
func UpdateSetValue(serverName, key, value, newValue string) error {
	err := RemoveSetValue(serverName, key, value)
	if err != nil {
		return err
	}

	err = AddValueToSet(serverName, key, newValue)
	if err != nil {
		return err
	}
	return nil
}

//RemoveSetValue removes a set member
func RemoveSetValue(serverName, key, value string) error {
	conn, err := connector.GetByName(serverName)
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
