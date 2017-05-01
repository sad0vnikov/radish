package db

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//ZSetKey is a Redis ZSET key
type ZSetKey struct {
	serverName string
	key        string
}

//ZSetMember is a ZSetMember struct
type ZSetMember struct {
	Score  int64
	Member string
}

//KeyType returns Zset key type
func (key ZSetKey) KeyType() string {
	return RedisZset
}

//PagesCount returns ZSET key values pages count
func (key ZSetKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		return 0, err
	}

	r, err := conn.Do("ZCARD", key.key)
	count, err := redis.Int(r, err)
	if err != nil {
		panic(err)
	}

	return getValuesPagesCount(count, pageSize), nil
}

//Values returns ZSET values page
func (key ZSetKey) Values(pageNum int, pageSize int) (interface{}, error) {
	conn, err := connector.GetByName(key.serverName)
	if err != nil {
		return nil, err
	}

	var (
		cursor int64
		values []string
	)

	r, err := redis.Values(conn.Do("ZSCAN", key.key, pageNum*pageSize, "COUNT", pageSize))
	r, err = redis.Scan(r, &cursor, &values)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var zSetValues []ZSetMember
	for i := 1; i < len(values); i = i + 2 {
		zsetMember := values[i-1]
		zsetScore, err := strconv.ParseInt(values[i], 0, 0)
		if err != nil {
			logger.Errorf("can't get convert %v score %v to string", zsetMember, zsetScore)
			return nil, err
		}

		zSetValues = append(zSetValues, ZSetMember{
			Score:  zsetScore,
			Member: zsetMember,
		})
	}

	return zSetValues, nil

}

//AddZSetValue adds a new sorted set value if it doesn't exist
func AddZSetValue(serverName, key, value string, score int64) error {
	conn, err := connector.GetByName(serverName)
	if err != nil {
		return err
	}

	_, err = conn.Do("ZADD", key, "NX", score, value)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//UpdateZSetValueIfExists updates a ZSet value if it exists
func UpdateZSetValueIfExists(serverName, key, oldValue, value string, score int64) error {
	err := DeleteZSetValue(serverName, key, oldValue)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = AddZSetValue(serverName, key, value, score)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//DeleteZSetValue deletes a ZSET member
func DeleteZSetValue(serverName, key, value string) error {
	conn, err := connector.GetByName(serverName)
	if err != nil {
		return err
	}

	_, err = conn.Do("ZREM", key, value)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
