package db

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/garyburd/redigo/redis"
	"github.com/sad0vnikov/radish/logger"
)

//Keys is a basic interface for managing redis keys
type Keys interface {
	getKeysByMask()
}

//ObjectType is a Redis object type
type ObjectType string

const (
	//RedisString is a Redis String type
	RedisString = "string"
	//RedisList is a Redis List type
	RedisList = "list"
	//RedisHash is a Redis Hash type
	RedisHash = "hash"
	//RedisSet is a Redis Set type
	RedisSet = "set"
	//RedisZset is a Redis Zset type
	RedisZset = "zset"
)

//Key is a redis key representation for Radish
type Key interface {
	Values(pageNum, pageSize int) (interface{}, error)
	PagesCount(pageSize int) (int, error)
	KeyType() string
}

var connector Connections

func init() {
	connector = &RedisConnections{}
}

//FindKeysByMask returns a list of keys satisfyig mask
func FindKeysByMask(serverName, mask string) ([]string, error) {

	conn, err := connector.GetByName(serverName)

	if err != nil {
		return nil, err
	}

	result, err := conn.Do("KEYS", mask)
	if err != nil {
		return nil, err
	}

	return redis.Strings(result, err)

}

//KeyExists returns True if given Redis key exists
func KeyExists(serverName, key string) (bool, error) {
	conn, err := connector.GetByName(serverName)

	if err != nil {
		logger.Error(err)
		return false, err
	}

	r, err := conn.Do("EXISTS", key)
	keyExists, err := redis.Bool(r, err)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	return keyExists, nil

}

//GetKeyInfo returns given key info
func GetKeyInfo(serverName, key string) (Key, error) {

	conn, err := connector.GetByName(serverName)

	if err != nil {
		log.Fatal(err)
		return nil, errors.New("can't connect to server " + serverName)
	}

	result, err := conn.Do("TYPE", key)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("can't get key %v type", key)
	}

	keyType, err := redis.String(result, err)
	switch keyType {
	case "string":
		return StringKey{serverName: serverName, key: key}, nil
	case "list":
		return ListKey{serverName: serverName, key: key}, nil
	case "hash":
		return HashKey{serverName: serverName, key: key}, nil
	case "set":
		return SetKey{serverName: serverName, key: key}, nil
	case "zset":
		return ZSetKey{serverName: serverName, key: key}, nil
	}

	return nil, errors.New("get unknown redis object type")
}

func getValuesPagesCount(valuesCout int, pageSize int) int {
	return int(math.Ceil(float64(valuesCout) / float64(pageSize)))
}
