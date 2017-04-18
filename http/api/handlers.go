package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/http/responds"
	"github.com/sad0vnikov/radish/http/server"
	"github.com/sad0vnikov/radish/logger"
	"github.com/sad0vnikov/radish/redis/db"
)

//GetServersList is a http handler returning a list of avalable Redis instances
func GetServersList(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return config.Get().Servers, nil
}

type getKeysByMaskResponse struct {
	Keys       []string
	Page       int
	PagesCount int
}

const defaultPageSize = 100

//GetKeysByMask is a http handler returning a JSON list of keys satisfying given mask
//for server with the name given in 'server' query param
func GetKeysByMask(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	const pageSize = defaultPageSize

	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is mandatory")
	}
	mask := requestParams["mask"]
	if len(mask) == 0 {
		return nil, responds.NewBadRequestError("'mask' param is mandatory")
	}

	pageNumber := 1
	page := requestParams["page"]
	if len(page) > 0 {
		paramPage, err := strconv.ParseInt(page, 0, 8)
		if err == nil {
			pageNumber = int(paramPage)
		}

	}

	keys, err := db.FindKeysByMask(serverName, mask)
	if err != nil {
		return nil, err
	}

	pageOffsetEnd := pageNumber * pageSize
	if pageOffsetEnd > len(keys) {
		pageOffsetEnd = len(keys)
	}

	pageOffsetStart := (pageNumber - 1) * pageSize
	if pageOffsetStart > len(keys) {
		return nil, responds.NewNotFoundError("page not found")
	}

	keysPage := keys[pageOffsetStart:pageOffsetEnd]
	pagesCount := int(math.Ceil(float64(len(keys)) / float64(pageSize)))

	responseContents := getKeysByMaskResponse{Keys: keysPage, Page: pageNumber, PagesCount: pagesCount}

	return responseContents, nil
}

type keyInfoResponse struct {
	PageSize   int
	PagesCount int
	KeyType    string
}

//GetKeyInfo returns key type, values pages count and page size
func GetKeyInfo(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is mandatory")
	}

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		return nil, responds.NewBadRequestError("'key' param is mandatory")
	}

	keyExists, err := db.KeyExists(serverName, keyName)
	if err != nil {
		return nil, err
	}
	if !keyExists {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	key, err := db.GetKeyInfo(serverName, keyName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	response := keyInfoResponse{}
	response.PageSize = defaultPageSize
	response.PagesCount, err = key.PagesCount(defaultPageSize)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	response.KeyType = key.KeyType()

	return response, nil
}

type valuesResponse struct {
	KeyType string
	Values  []RedisValue
	PageNum int
}

//RedisValue is a struct that could represent any Redis value
//Value field is always present
//Key fields is only present in a redis Hash values
//Score field is only present in redis ZSet values
type RedisValue struct {
	Value string
	Key   string
	Score int64
}

//GetKeyValues returns a list of key values
func GetKeyValues(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is required")
	}

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		return nil, responds.NewBadRequestError("'key' param is required")
	}

	keyExists, err := db.KeyExists(serverName, keyName)
	if err != nil {
		return nil, err
	}
	if !keyExists {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	pageParam := requestParams["page"]
	pageNum := 0
	if len(pageParam) != 0 {
		if parsedPageParam, err := strconv.ParseInt(pageParam, 0, 0); err == nil {
			pageNum = int(parsedPageParam)
		}
	}

	key, err := db.GetKeyInfo(serverName, keyName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	v, err := key.Values(pageNum, defaultPageSize)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	response := valuesResponse{}
	response.PageNum = pageNum
	switch key.KeyType() {
	case db.RedisString:
		response.KeyType = key.KeyType()
		if str, ok := v.(string); ok {
			value := RedisValue{Value: str}
			response.Values = []RedisValue{value}
		}
		return response, nil

	case db.RedisList:
		response.KeyType = key.KeyType()
		if strings, ok := v.([]string); ok {
			values := []RedisValue{}
			for _, str := range strings {
				values = append(values, RedisValue{Value: str})
			}
			response.Values = values
		}
		return response, nil

	case db.RedisZset:
		response.KeyType = key.KeyType()
		if v, ok := v.([]db.ZSetMember); ok {
			values := []RedisValue{}
			for _, setMember := range v {
				values = append(values, RedisValue{Value: setMember.Member, Score: setMember.Score})
			}
			response.Values = values
		}
		response.PageNum = pageNum
		return response, nil
	case db.RedisHash:
		response.KeyType = key.KeyType()
		values := []RedisValue{}
		if v, ok := v.(map[string]string); ok {
			for hashKey, hashValue := range v {
				values = append(values, RedisValue{Key: hashKey, Value: hashValue})
			}
		}
		response.Values = values
		return response, nil
	case db.RedisSet:
		response.KeyType = key.KeyType()
		values := []RedisValue{}
		if v, ok := v.([]string); ok {
			for _, setMember := range v {
				values = append(values, RedisValue{Value: setMember})
			}
		}
		response.Values = values
		return response, nil

	default:
		return nil, fmt.Errorf("%v key has not-supported type %v", keyName, key.KeyType())
	}

}
