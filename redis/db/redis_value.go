package db

import (
	"unicode"
)

//RedisValue represents any redis key's value
type RedisValue struct {
	Value    string
	IsBinary bool
}

func isBinary(s string) bool {
	for _, ch := range s {
		if ch > unicode.MaxASCII || !unicode.IsPrint(ch) {
			return true
		}
	}

	return false
}
