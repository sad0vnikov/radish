package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//ListKey is a key for redis List
type ListKey struct {
	key        string
	serverName string
}

//KeyType returns List key type
func (key ListKey) KeyType() string {
	return RedisList
}

//PagesCount returns List values pages count
func (key ListKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		return 0, err
	}
	r, err := conn.Do("LLEN", key.key)
	count, err := redis.Int(r, err)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return getValuesPagesCount(count, pageSize), nil
}

//Values returns redis List values page
func (key ListKey) Values(pageNum, pageSize int) (interface{}, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		return 0, err
	}

	pageStart := pageNum * pageSize
	pageEnd := (pageNum + 1) * pageSize
	r, err := conn.Do("LRANGE", key.key, pageStart, pageEnd)
	values, err := redis.Strings(r, err)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return values, nil
}
