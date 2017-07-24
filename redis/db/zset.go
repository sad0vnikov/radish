package db

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
	rd "github.com/sad0vnikov/radish/redis"
)

type ZSetValues struct {
	values           []ZSetMember
	pagesCount       int
	valuesLoaded     bool
	pagesCountLoaded bool
	query            *KeyValuesQuery
	key              ZSetKey
}

func (v *ZSetValues) Values() (interface{}, error) {
	if !v.valuesLoaded {
		loadedValues, err := v.key.getValues(v.query.PageNum, v.query.PageSize)
		if err != nil {
			return nil, err
		}
		v.values = loadedValues
	}
	return v.values, nil
}

func (v *ZSetValues) PagesCount() (int, error) {
	if !v.pagesCountLoaded {
		pagesCount, err := v.key.getPagesCount(v.query.PageSize)
		if err != nil {
			return 0, err
		}
		v.pagesCount = pagesCount
	}
	return v.pagesCount, nil
}

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

func (key ZSetKey) Values(query *KeyValuesQuery) KeyValues {
	return &ZSetValues{
		key:   key,
		query: query,
	}

}

//PagesCount returns ZSET key values pages count
func (key ZSetKey) getPagesCount(pageSize int) (int, error) {
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
func (key ZSetKey) getValues(pageNum int, pageSize int) ([]ZSetMember, error) {
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
