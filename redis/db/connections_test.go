package db

import (
	"reflect"
	"testing"
	"github.com/sad0vnikov/radish/redis"
)

func TestCorrectStringParses(t *testing.T) {
	s := "db0:keys=5,expires=0,avg_ttl=0"
	expected := redis.ServerKeyspaceStat{KeysCount: 5}

	result := parseKeyspaceStatString(s)
	eq := reflect.DeepEqual(result, expected)
	if !eq {
		t.Errorf("got wrong parsed string '%s': %v, expected: %v", s, result, expected)
	}
}
