package db

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
)

//ZSetKey is a Redis ZSET key
type ZSetKey struct {
	serverName string
	dbNum      uint8
	key        string
}

//ZSetMember is a ZSetMember struct
type ZSetMember struct {
	Score  int64
	Member RedisValue
}

//KeyType returns Zset key type
func (key ZSetKey) KeyType() string {
	return RedisZset
}

//PagesCount returns ZSET key values pages count
func (key ZSetKey) PagesCount(pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName, key.dbNum)
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
	conn, err := connector.GetByName(key.serverName, key.dbNum)
	if err != nil {
		return nil, err
	}

	var (
		values []string
	)

	r, err := conn.Do("ZRANGEBYSCORE", key.key, "-inf", "+inf", "WITHSCORES")
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

	var zSetValues []ZSetMember
	for i := 1; i < len(valuesPage); i = i + 2 {
		zsetMember := valuesPage[i-1]
		zsetScore, err := strconv.ParseInt(valuesPage[i], 0, 0)
		if err != nil {
			logger.Errorf("can't get convert %v score %v to string", zsetMember, zsetScore)
			return nil, err
		}

		zSetValues = append(zSetValues, ZSetMember{
			Score: zsetScore,
			Member: RedisValue{
				Value:    zsetMember,
				IsBinary: isBinary(zsetMember),
			},
		})
	}

	return zSetValues, nil

}

//AddZSetValue adds a new sorted set value if it doesn't exist
func AddZSetValue(serverName string, dbNum uint8, key, value string, score int64) error {
	conn, err := connector.GetByName(serverName, dbNum)
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
func UpdateZSetValueIfExists(serverName string, dbNum uint8, key, oldValue, value string, score int64) error {
	err := DeleteZSetValue(serverName, dbNum, key, oldValue)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = AddZSetValue(serverName, dbNum, key, value, score)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//DeleteZSetValue deletes a ZSET member
func DeleteZSetValue(serverName string, dbNum uint8, key, value string) error {
	conn, err := connector.GetByName(serverName, dbNum)
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
