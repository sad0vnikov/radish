package db

import (
	"math"

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
		logger.Error(err)
		return 0, err
	}
	r, err := conn.Do("LLEN", key.key)
	count, err := redis.Int(r, err)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return int(math.Ceil(float64(count) / float64(pageSize))), nil
}

//Values returns redis List values page
func (key ListKey) Values(pageNum, pageSize int) (interface{}, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	pageStart := pageNum * pageSize
	pageEnd := (pageNum + 1) * pageSize
	r, err := conn.Do("LRANGE", key.key, pageStart, pageEnd)
	values, err := redis.Strings(r, err)

	if err != nil {
		logger.Error(err)
		return values, err
	}

	return values, nil
}
