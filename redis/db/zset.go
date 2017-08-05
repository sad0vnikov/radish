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
	totalValuesCount int
	query            *KeyValuesQuery
	key              ZSetKey
}

func (v *ZSetValues) Values() (interface{}, error) {
	if !v.valuesLoaded {
		err := v.loadValues()
		if err != nil {
			return nil, err
		}
	}
	return v.values, nil
}

func (v *ZSetValues) PagesCount() (int, error) {
	if !v.pagesCountLoaded {
		pagesCount, err := v.calculatePagesCount()
		if err != nil {
			return 0, err
		}
		v.pagesCount = pagesCount
	}
	return v.pagesCount, nil
}

func (vInfo *ZSetValues) TotalValuesCount() (int, error) {
	var err error
	if !vInfo.valuesLoaded {
		err = vInfo.loadValues()
	}
	return vInfo.totalValuesCount, err
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

//PagesCount returns ZSET key vInfo pages count
func (vInfo *ZSetValues) calculatePagesCount() (int, error) {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return 0, err
	}

	count := 0
	if vInfo.query.Mask == "*" {
		r, err := conn.Do("ZCARD", vInfo.key.key)
		count, err = redis.Int(r, err)
	} else {
		if !vInfo.valuesLoaded {
			err = vInfo.loadValues()
		}
		count = vInfo.totalValuesCount
	}

	if err != nil {
		panic(err)
	}

	return getValuesPagesCount(count, vInfo.query.PageSize), nil
}

//Values returns ZSET vInfo page
func (vInfo *ZSetValues) loadValues() error {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return err
	}

	var (
		values []string
	)

	r, err := conn.Do("ZRANGEBYSCORE", vInfo.key.key, "-inf", "+inf", "WITHSCORES")
	values, err = redis.Strings(r, err)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err != nil {
		return err
	}

	var zSetValues []ZSetMember
	maskedValuesCount := 0
	offsetStart, offsetEnd, err := rd.GetPageRangeForStrings(values, vInfo.query.PageSize*2, vInfo.query.PageNum)
	if err != nil {
		return err
	}
	for i := 1; i < len(values); i = i + 2 {
		zsetMember := values[i-1]
		zsetScore, err := strconv.ParseInt(values[i], 0, 0)
		if err != nil {
			logger.Errorf("can't get convert %vInfo score %vInfo to string", zsetMember, zsetScore)
			return err
		}

		matchesMask := matchStringValueWithMask(zsetMember, vInfo.query.Mask)
		if matchesMask && i >= offsetStart && i <= offsetEnd {
			zSetValues = append(zSetValues, ZSetMember{
				Score: zsetScore,
				Member: RedisValue{
					Value:    zsetMember,
					IsBinary: isBinary(zsetMember),
				},
			})
		}

		if matchesMask {
			maskedValuesCount++
		}

	}

	vInfo.totalValuesCount = maskedValuesCount

	vInfo.values = zSetValues

	return nil

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
