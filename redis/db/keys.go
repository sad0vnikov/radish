package db

import (
	"github.com/garyburd/redigo/redis"
)

//Keys is a basic interface for managing redis keys
type Keys interface {
	getKeysByMask()
}

var connector Connections

func init() {
	connector = &RedisConnections{}
}

//FindKeysByMask returns a list of keys satisfyig mask
func FindKeysByMask(serverName string, mask string) ([]string, error) {

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
