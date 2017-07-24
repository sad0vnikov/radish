package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//ListValues stores List vInfo data
type ListValues struct {
	values           []ListMember
	pagesCount       int
	pagesCountLoaded bool
	valuesLoaded     bool
	query            *KeyValuesQuery
	key              ListKey
}

//Values returns redis List vInfo page
func (vInfo *ListValues) Values() (interface{}, error) {
	if !vInfo.valuesLoaded {
		err := vInfo.loadValues()
		if err != nil {
			return nil, err
		}
	}

	return vInfo.values, nil
}

//PagesCount returns List vInfo pages count
func (vInfo *ListValues) PagesCount() (int, error) {
	if !vInfo.pagesCountLoaded {
		pagesCount, err := vInfo.calculatePagesCount()
		if err != nil {
			return 0, err
		}
		vInfo.pagesCount = pagesCount
	}

	return vInfo.pagesCount, nil
}

//ListKey is a key for redis List
type ListKey struct {
	key        string
	serverName string
	dbNum      uint8
}

//Values returns ListValues object
func (key ListKey) Values(query *KeyValuesQuery) KeyValues {
	return &ListValues{
		key:   key,
		query: query,
	}
}

//KeyType returns List key type
func (key ListKey) KeyType() string {
	return RedisList
}

//ListMember represents List member value
type ListMember struct {
	Index int
	Value RedisValue
}

func (vInfo *ListValues) calculatePagesCount() (int, error) {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return 0, err
	}
	count := 0
	if vInfo.query.Mask == "*" {
		r, err := conn.Do("LLEN", vInfo.key.key)
		count, err = redis.Int(r, err)

	} else {
		if !vInfo.pagesCountLoaded {
			err = vInfo.loadValues()
		}
		count = len(vInfo.values)
	}

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return getValuesPagesCount(count, vInfo.query.PageSize), nil
}

//Values returns redis List vInfo page
func (vInfo *ListValues) loadValues() error {
	conn, err := connector.GetByName(vInfo.key.serverName, vInfo.key.dbNum)
	if err != nil {
		return err
	}

	pageStart := (vInfo.query.PageNum - 1) * vInfo.query.PageSize
	pageEnd := vInfo.query.PageNum*vInfo.query.PageSize - 1
	r, err := conn.Do("LRANGE", vInfo.key.key, pageStart, pageEnd)
	strings, err := redis.Strings(r, err)

	if err != nil {
		logger.Error(err)
		return nil
	}

	var values = make([]ListMember, len(strings))
	memberIndex := pageStart
	i := 0
	for _, v := range strings {
		if matchStringValueWithMask(v, vInfo.query.Mask) {
			rv := RedisValue{Value: v, IsBinary: isBinary(v)}
			values[i] = ListMember{Value: rv, Index: memberIndex}
			i++
		}
		memberIndex++
	}

	vInfo.values = values[:i]
	return nil
}

//InsertToListWithPos inserts a value at the given position
//If there are vInfo after the given index, they are moved to the right
//If position greater then the last list index, the value will be added to the and of the list
func InsertToListWithPos(serverName string, dbNum uint8, key, listValue string, position int) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	r, err := conn.Do("LINDEX", key, position)
	valueAfter, err := redis.String(r, err)
	if err != nil && err != redis.ErrNil {
		logger.Error(err)
		return err
	}

	if len(valueAfter) != 0 {
		_, err = conn.Do("LINSERT", key, "BEFORE", valueAfter, listValue)
	} else {
		_, err = conn.Do("RPUSH", key, listValue)
	}

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//AppendToList appends a value to the end of list
func AppendToList(serverName string, dbNum uint8, key, listValue string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", key, listValue)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//UpdateListValue updates a list Value by index
func UpdateListValue(serverName string, dbNum uint8, key string, index int, newValue string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	_, err = conn.Do("LSET", key, index, newValue)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//DeleteListValue removes List member
func DeleteListValue(serverName string, dbNum uint8, key string, index int) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		return err
	}

	const deletedValue = "RADISH_DELETED"
	_, err = conn.Do("LSET", key, index, deletedValue)
	if err != nil {
		logger.Error(err)
		return err
	}

	_, err = conn.Do("LREM", key, 1, deletedValue)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
