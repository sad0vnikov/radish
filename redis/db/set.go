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

//SetValues represents set vInfo
type SetValues struct {
	values          []RedisValue
	pagesCount      int
	valuesLoaded    bool
	pageCountLoaded bool
	query           *KeyValuesQuery
	key             SetKey
}

//Values returns vInfo for a set
func (vInfo *SetValues) Values() (interface{}, error) {
	if !vInfo.valuesLoaded {
		err := vInfo.loadSetValues()
		if err != nil {
			return nil, err
		}
	}
	return vInfo.values, nil
}

//PagesCount returns pagesCount for set vInfo
func (vInfo *SetValues) PagesCount() (int, error) {
	if !vInfo.pageCountLoaded {
		pagesCount, err := vInfo.calculatePagesCount()
		if err != nil {
			return 0, err
		}
		vInfo.pagesCount = pagesCount
	}
	return vInfo.pagesCount, nil
}

//KeyType returns Redis SET type
func (key SetKey) KeyType() string {
	return RedisSet
}

//Values returns SetValues object
func (key SetKey) Values(query *KeyValuesQuery) KeyValues {

	return &SetValues{
		key:   key,
		query: query,
	}
}

//PagesCount returns Redis SET vInfo pages count
func (vInfo *SetValues) calculatePagesCount() (int, error) {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return 0, err
	}

	count := 0
	if vInfo.query.Mask == "*" {
		r, err := conn.Do("SCARD", vInfo.key.key)
		count, err = redis.Int(r, err)
	} else {
		if !vInfo.valuesLoaded {
			err = vInfo.loadSetValues()
		}
		count = len(vInfo.values)
	}

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return getValuesPagesCount(count, vInfo.query.PageSize), nil
}

//Values returns a SET vInfo page
func (vInfo *SetValues) loadSetValues() error {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		logger.Error(err)
		return err
	}

	var (
		values []string
	)

	r, err := conn.Do("SMEMBERS", vInfo.key.key)
	values, err = redis.Strings(r, err)
	if err != nil {
		logger.Error(err)
		return err
	}
	offsetStart, offsetEnd, err := rd.GetPageRangeForStrings(values, vInfo.query.PageSize, vInfo.query.PageNum)

	if err != nil {
		return err
	}
	valuesPageStrings := values[offsetStart:offsetEnd]

	valuesPage := make([]RedisValue, len(valuesPageStrings))
	i := 0
	for _, s := range valuesPageStrings {
		if matchStringValueWithMask(s, vInfo.query.Mask) {
			valuesPage[i] = RedisValue{
				Value:    s,
				IsBinary: isBinary(s),
			}
			i++
		}
	}

	vInfo.values = valuesPage[:i]
	return nil
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
