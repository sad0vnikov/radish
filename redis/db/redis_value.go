package db

import (
	"unicode"
)

const defaultPageSize = 100

//KeyValues is an interface representing key's values
type KeyValues interface {
	Values() (interface{}, error)
	PagesCount() (int, error)
}

type KeyValuesQuery struct {
	PageNum  int
	PageSize int
	Mask     string
}

func NewKeyValuesQuery() *KeyValuesQuery {
	return &KeyValuesQuery{PageNum: 1, PageSize: defaultPageSize, Mask: "*"}
}

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
