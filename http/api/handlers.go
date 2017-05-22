package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

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
	mask := r.URL.Query().Get("mask")
	if len(mask) == 0 {
		mask = "*"
	}

	pageNumber := 1
	page := r.URL.Query().Get("page")
	if len(page) > 0 {
		paramPage, err := strconv.ParseInt(page, 0, 8)
		if err == nil {
			pageNumber = int(paramPage)
		}

	}

	keys, err := db.FindKeysByMask(serverName, mask)
	if err != nil {
		logger.Error(err)
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

//KeysSubtreeResponse is a
type KeysSubtreeResponse struct {
	Nodes     []db.KeyTreeNode
	Path      []string
	KeysCount int64
}

//GetKeysSubtree is a handler for getting Redis keys tree nodes
func GetKeysSubtree(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is mandatory")
	}

	keyPrefix := requestParams["prefix"]
	if len(keyPrefix) == 0 {
		return nil, responds.NewBadRequestError("'prefix' param is mandatory")
	}

	offsetParam := requestParams["offset"]
	offset, err := strconv.ParseInt(offsetParam, 0, 0)
	if err != nil {
		offset = 0
	}

	delimiter := ":"
	var pageSize int32
	pageSize = 100

	node := db.KeyTreeNode{Name: keyPrefix, HasChildren: true}
	nodes, keysCount, err := db.FindKeysTreeNodeChildren(serverName, delimiter, offset, pageSize, node)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	path := []string{}
	if keyPrefix != "*" {
		path = strings.Split(keyPrefix, delimiter)
	}
	response := KeysSubtreeResponse{Nodes: nodes, KeysCount: keysCount, Path: path}

	return response, nil
}

type singleValueResponse struct {
	KeyType string
	Value   string
}

type listValuesResponse struct {
	KeyType    string
	Values     []db.ListMember
	PageNum    int
	PagesCount int
}

type hashValuesResponse struct {
	KeyType    string
	Values     map[string]string
	PageNum    int
	PagesCount int
}

type setValuesResponse struct {
	KeyType    string
	Values     []string
	PageNum    int
	PagesCount int
}

type zsetValuesResponse struct {
	KeyType    string
	Values     []db.ZSetMember
	PageNum    int
	PagesCount int
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

	pagesCount, err := key.PagesCount(defaultPageSize)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	switch key.KeyType() {
	case db.RedisString:
		response := singleValueResponse{}
		response.KeyType = key.KeyType()
		if str, ok := v.(string); ok {
			response.Value = str
		}
		return response, nil

	case db.RedisList:
		response := listValuesResponse{}
		response.KeyType = key.KeyType()
		if strings, ok := v.([]db.ListMember); ok {
			response.Values = strings
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil

	case db.RedisZset:
		response := zsetValuesResponse{}
		response.KeyType = key.KeyType()
		if v, ok := v.([]db.ZSetMember); ok {
			response.Values = v
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil
	case db.RedisHash:
		response := hashValuesResponse{}
		response.KeyType = key.KeyType()
		if v, ok := v.(map[string]string); ok {
			response.Values = v
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil
	case db.RedisSet:
		response := setValuesResponse{}
		response.KeyType = key.KeyType()
		if v, ok := v.([]string); ok {
			response.Values = v
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil

	default:
		return nil, fmt.Errorf("%v key has not-supported type %v", keyName, key.KeyType())
	}

}

//DeleteKey deletes a given key
func DeleteKey(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	requestParams := server.GetURLParams(r)
	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is required")
	}

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		return nil, responds.NewBadRequestError("'key' param is required")
	}

	err := db.DeleteKey(serverName, keyName)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type addStringJSONRequest struct {
	Value string
}

//AddStringValue adds a new string value
func AddStringValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, responds.NewBadRequestError(err.Error())
	}
	serverName := GetParam("server", r)
	keyName := GetParam("key", r)

	decoder := json.NewDecoder(r.Body)
	var bodyReq addStringJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}

	ex, err := db.KeyExists(serverName, keyName)
	if err != nil {
		return nil, err
	}

	if ex {
		return nil, responds.NewConflictError(fmt.Sprintf("key %v already exists", keyName))
	}

	err = db.Set(serverName, keyName, bodyReq.Value)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type updateStringRequest struct {
	Value string
}

//UpdateStringValue updates a string value
func UpdateStringValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(r.Body)
	var JSONReq updateStringRequest
	err = decoder.Decode(&JSONReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}

	serverName := GetParam("server", r)
	keyName := GetParam("key", r)

	ex, err := db.KeyExists(serverName, keyName)
	if err != nil {
		return nil, err
	}

	if !ex {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}
	if len(JSONReq.Value) == 0 {
		return nil, responds.NewBadRequestError("`Value` JSON param is missing")
	}
	err = db.Set(serverName, keyName, JSONReq.Value)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type addHashValueJSONRequest struct {
	Key   string
	Value string
}

//AddHashValue adds a new redis Hash value
func AddHashValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, responds.NewBadRequestError(err.Error())
	}
	serverName := GetParam("server", r)
	keyName := GetParam("key", r)

	decoder := json.NewDecoder(r.Body)
	var bodyReq addHashValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Key) == 0 {
		return nil, responds.NewBadRequestError("JSON `Key` param is missing")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}

	hashKey := bodyReq.Key
	hashValue := bodyReq.Value

	ex, err := db.HashKeyExists(serverName, keyName, hashKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if ex {
		return nil, responds.NewConflictError(fmt.Sprintf("key %v already exists", keyName))
	}

	err = db.SetHashKey(serverName, keyName, hashKey, hashValue)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type updateHashValueJSONRequest struct {
	Value string
}

//UpdateHashValue updates an exists hash value
func UpdateHashValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "hashKey"}, r)
	if err != nil {
		return nil, responds.NewBadRequestError(err.Error())
	}
	serverName := GetParam("server", r)
	keyName := GetParam("key", r)
	hashKey := GetParam("hashKey", r)

	decoder := json.NewDecoder(r.Body)
	var bodyReq updateHashValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}
	hashValue := bodyReq.Value

	ex, err := db.HashKeyExists(serverName, keyName, hashKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !ex {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	err = db.SetHashKey(serverName, keyName, hashKey, hashValue)
	if err != nil {
		return nil, err
	}

	return "", nil
}

//DeleteHashValue updates an exists hash value
func DeleteHashValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "hashKey"}, r)
	if err != nil {
		return nil, responds.NewBadRequestError(err.Error())
	}
	serverName := GetParam("server", r)
	keyName := GetParam("key", r)
	hashKey := GetParam("hashKey", r)

	ex, err := db.HashKeyExists(serverName, keyName, hashKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !ex {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	err = db.DeleteHashValue(serverName, keyName, hashKey)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type addToListJSONRequest struct {
	Value string
	Index *int
}

//AddListValue adds a new List value
func AddListValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var bodyReq addToListJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}

	if bodyReq.Index != nil {
		err = db.InsertToListWithPos(GetParam("server", r), GetParam("key", r), bodyReq.Value, *bodyReq.Index)
	} else {
		err = db.AppendToList(GetParam("server", r), GetParam("key", r), bodyReq.Value)
	}

	if err != nil {
		return nil, err
	}

	return "", nil
}

type updateListValueJSONRequest struct {
	Value string
}

//UpdateListValue updates a list value
func UpdateListValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "index"}, r)
	if err != nil {
		return nil, err
	}

	indexParam := GetParam("index", r)
	index, err := strconv.ParseInt(indexParam, 0, 0)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if index > math.MaxInt32 {
		logger.Error(err)
		return nil, responds.NewBadRequestError("'index' is too large")
	}

	decoder := json.NewDecoder(r.Body)
	var bodyReq updateListValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}

	err = db.UpdateListValue(GetParam("server", r), GetParam("key", r), int(index), bodyReq.Value)
	if err != nil {
		return nil, err
	}

	return "", nil
}

//DeleteListValue updates a list value
func DeleteListValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "index"}, r)
	if err != nil {
		return nil, err
	}

	indexParam := GetParam("index", r)
	index, err := strconv.ParseInt(indexParam, 0, 0)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if index > math.MaxInt32 {
		logger.Error(err)
		return nil, responds.NewBadRequestError("'index' is too large")
	}

	err = db.DeleteListValue(GetParam("server", r), GetParam("key", r), int(index))
	if err != nil {
		return nil, err
	}

	return "", nil
}

type addSetValueJSONRequest struct {
	Value string
}

//AddSetValue adds a new SET member
func AddSetValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var bodyReq addSetValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}
	err = db.AddValueToSet(GetParam("server", r), GetParam("key", r), bodyReq.Value)
	if err != nil {
		return "", err
	}

	return "", err
}

type updateSetValueJSONRequest struct {
	Value string
}

//UpdateSetValue updates a set member
func UpdateSetValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var bodyReq updateSetValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}

	err = db.UpdateSetValue(GetParam("server", r), GetParam("key", r), GetParam("value", r), bodyReq.Value)
	if err != nil {
		return "", err
	}

	return "", err
}

//DeleteSetValue updates a set member
func DeleteSetValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "value"}, r)
	if err != nil {
		return nil, err
	}

	err = db.DeleteSetValue(GetParam("server", r), GetParam("key", r), GetParam("value", r))
	if err != nil {
		return "", err
	}

	return "", err
}

type addZSetValueJSONRequest struct {
	Value string
	Score int64
}

//AddZSetValue adds a new ZSET value
func AddZSetValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key"}, r)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var bodyReq addZSetValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}

	err = db.AddZSetValue(GetParam("server", r), GetParam("key", r), bodyReq.Value, bodyReq.Score)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type updateZSetValueJSONRequest struct {
	Value string
	Score int64
}

//UpdateZSetValue updates a ZSET value
func UpdateZSetValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "value"}, r)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r.Body)
	var bodyReq updateZSetValueJSONRequest
	err = decoder.Decode(&bodyReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}
	if len(bodyReq.Value) == 0 {
		return nil, responds.NewBadRequestError("JSON `Value` param is missing")
	}
	err = db.UpdateZSetValueIfExists(GetParam("server", r), GetParam("key", r), GetParam("value", r), bodyReq.Value, bodyReq.Score)
	if err != nil {
		return nil, err
	}

	return "", nil
}

//DeleteZSetValue updates a ZSET value
func DeleteZSetValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "value"}, r)
	if err != nil {
		return nil, err
	}

	err = db.DeleteZSetValue(GetParam("server", r), GetParam("key", r), GetParam("value", r))
	if err != nil {
		return nil, err
	}

	return "", nil
}
