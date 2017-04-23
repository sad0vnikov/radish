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

//InsertToList inserts a value at the given position
//If there are values after the given index, they are moved to the right
//If position greater then the last list index, the value will be added to the and of the list
func InsertToList(serverName, key, listValue string, position int) error {
	conn, err := connector.GetByName(serverName)
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

//UpdateListValue updates a list Value by index
func UpdateListValue(serverName, key string, index int, newValue string) error {
	conn, err := connector.GetByName(serverName)
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
func DeleteListValue(serverName, key string, index int) error {
	conn, err := connector.GetByName(serverName)
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
