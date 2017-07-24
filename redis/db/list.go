package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//ListValues stores List values data
type ListValues struct {
	values           []ListMember
	pagesCount       int
	pagesCountLoaded bool
	valuesLoaded     bool
	query            *KeyValuesQuery
	key              ListKey
}

//Values returns redis List values page
func (values *ListValues) Values() (interface{}, error) {
	if !values.valuesLoaded {
		loadedValues, err := values.key.getValuesForString(values.query.Mask, values.query.PageNum, values.query.PageSize)
		if err != nil {
			return nil, err
		}
		values.values = loadedValues
	}

	return values.values, nil
}

//PagesCount returns List values pages count
func (values *ListValues) PagesCount() (int, error) {
	if !values.pagesCountLoaded {
		pagesCount, err := values.key.calculatePagesCount(values.query.Mask, values.query.PageSize)
		if err != nil {
			return 0, err
		}
		values.pagesCount = pagesCount
	}

	return values.pagesCount, nil
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

func (key *ListKey) calculatePagesCount(mask string, pageSize int) (int, error) {
	conn, err := connector.GetByName(key.serverName, key.dbNum)
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
func (key *ListKey) getValuesForString(mask string, pageNum, pageSize int) ([]ListMember, error) {
	conn, err := connector.GetByName(key.serverName, key.dbNum)
	if err != nil {
		return nil, err
	}

	pageStart := (pageNum - 1) * pageSize
	pageEnd := pageNum*pageSize - 1
	r, err := conn.Do("LRANGE", key.key, pageStart, pageEnd)
	strings, err := redis.Strings(r, err)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var values = make([]ListMember, len(strings))
	memberIndex := pageStart
	for i, v := range strings {
		rv := RedisValue{Value: v, IsBinary: isBinary(v)}
		values[i] = ListMember{Value: rv, Index: memberIndex}
		memberIndex++
	}

	return values, nil
}

//InsertToListWithPos inserts a value at the given position
//If there are values after the given index, they are moved to the right
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
