package db

import (
	"errors"
	"fmt"
	"math"
	"strings"

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
	Values(*KeyValuesQuery) KeyValues
	KeyType() string
}

//KeyTreeNode is a node of keys tree
type KeyTreeNode struct {
	Name        string
	Key         string
	HasChildren bool
}

var connector Connections

func init() {
	connector = &RedisConnections{}
}

//FindKeysByMask returns a list of keys satisfyig mask
func FindKeysByMask(serverName string, dbNum uint8, mask string) ([]string, error) {

	conn, err := connector.GetByName(serverName, dbNum)

	if err != nil {
		return nil, err
	}

	result, err := conn.Do("KEYS", mask)
	if err != nil {
		return nil, err
	}

	return redis.Strings(result, err)

}

//FindKeysTreeNodeChildren returns the first generation of children for given keys tree node
func FindKeysTreeNodeChildren(serverName string, dbNum uint8, delimiter string, pageSize int32, node KeyTreeNode) ([]KeyTreeNode, error) {
	var maskForSearch = node.Name
	if maskForSearch != "*" {
		maskForSearch = maskForSearch + delimiter + "*"
	}

	conn, err := connector.GetByName(serverName, dbNum)

	if err != nil {
		return nil, err
	}

	r, err := conn.Do("KEYS", maskForSearch)
	if err != nil {
		return nil, err
	}

	keys, _ := redis.Strings(r, nil)

	return getChildrenFromKeys(keys, maskForSearch, delimiter), nil
}

func getChildrenFromKeys(keys []string, maskForSearch, delimiter string) []KeyTreeNode {
	var childKeysMap = make(map[string]KeyTreeNode)

	var keyPrefix = ""
	if maskForSearch != "*" {
		keyPrefix = strings.TrimSuffix(maskForSearch, "*")
	}

	for _, key := range keys {
		nodeName := strings.TrimPrefix(key, keyPrefix)
		sepIndex := strings.Index(nodeName, delimiter)
		hasChildren := false
		if sepIndex != -1 {
			hasChildren = true
			nodeName = strings.Split(nodeName, delimiter)[0]
		}

		mapKey := nodeName
		if hasChildren {
			mapKey += delimiter
		}

		if _, prs := childKeysMap[mapKey]; prs == false {
			node := KeyTreeNode{Name: nodeName, HasChildren: hasChildren}
			if !hasChildren {
				node.Key = key
			}
			childKeysMap[mapKey] = node
		}
	}

	childKeys := make([]KeyTreeNode, len(childKeysMap))
	i := 0
	for _, k := range childKeysMap {
		childKeys[i] = k
		i++
	}

	return childKeys
}

//KeyExists returns True if given Redis key exists
func KeyExists(serverName string, dbNum uint8, key string) (bool, error) {
	conn, err := connector.GetByName(serverName, dbNum)

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
func GetKeyInfo(serverName string, dbNum uint8, key string) (Key, error) {

	conn, err := connector.GetByName(serverName, dbNum)

	if err != nil {
		logger.Critical(err)
		return nil, errors.New("can't connect to server " + serverName)
	}

	result, err := conn.Do("TYPE", key)
	if err != nil {
		logger.Critical(err)
		return nil, fmt.Errorf("can't get key %vInfo type", key)
	}

	keyType, err := redis.String(result, err)
	switch keyType {
	case "string":
		return StringKey{serverName: serverName, dbNum: dbNum, key: key}, nil
	case "list":
		return ListKey{serverName: serverName, dbNum: dbNum, key: key}, nil
	case "hash":
		return HashKey{serverName: serverName, dbNum: dbNum, key: key}, nil
	case "set":
		return SetKey{serverName: serverName, dbNum: dbNum, key: key}, nil
	case "zset":
		return ZSetKey{serverName: serverName, dbNum: dbNum, key: key}, nil
	}

	return nil, errors.New("get unknown redis object type")
}

//DeleteKey deletes a given key
func DeleteKey(serverName string, dbNum uint8, key string) error {
	conn, err := connector.GetByName(serverName, dbNum)
	if err != nil {
		logger.Critical(err)
		return err
	}

	_, err = conn.Do("DEL", key)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func getValuesPagesCount(valuesCout int, pageSize int) int {
	return int(math.Ceil(float64(valuesCout) / float64(pageSize)))
}
