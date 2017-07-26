package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"strings"

	"io/ioutil"

	"github.com/sad0vnikov/radish/http/responds"
	"github.com/sad0vnikov/radish/http/server"
	"github.com/sad0vnikov/radish/logger"
	"github.com/sad0vnikov/radish/redis"
	"github.com/sad0vnikov/radish/redis/db"
)

//GetServersList is a http handler returning a list of avalable Redis instances
func GetServersList(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return db.GetServersWithConnectionData(), nil
}

type getKeysByMaskResponse struct {
	Keys           []string
	FoundKeysCount int
	Page           int
	PagesCount     int
}

const defaultPageSize = 100

//GetMaxDbNumber is a handler used to get a maxium db number value for given server name
func GetMaxDbNumber(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server"}, r)
	if err != nil {
		return nil, responds.NewBadRequestError(err.Error())
	}

	serverName := GetParam("server", r)
	maxConnections, err := db.GetMaxDbNumsForServer(serverName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return int(maxConnections), nil
}

//GetKeysByMask is a http handler returning a JSON list of keys satisfying given mask
//for server with the name given in 'server' query param
func GetKeysByMask(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	const pageSize = defaultPageSize

	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is mandatory")
	}
	var dbNum uint8
	dbNum, err := GetParamUint8("db", r)

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

	keys, err := db.FindKeysByMask(serverName, dbNum, mask)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	pageOffsetStart, pageOffsetEnd, err := redis.GetPageRangeForStrings(keys, pageSize, pageNumber)
	if err != nil {
		return nil, responds.NewNotFoundError(err.Error())
	}

	keysPage := keys[pageOffsetStart:pageOffsetEnd]
	pagesCount := int(math.Ceil(float64(len(keys)) / float64(pageSize)))
	keysCount := len(keys)

	responseContents := getKeysByMaskResponse{Keys: keysPage, Page: pageNumber, PagesCount: pagesCount, FoundKeysCount: keysCount}

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
	var dbNum uint8
	dbNum, err := GetParamUint8("db", r)

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		return nil, responds.NewBadRequestError("'key' param is mandatory")
	}

	keyExists, err := db.KeyExists(serverName, dbNum, keyName)
	if err != nil {
		return nil, err
	}
	if !keyExists {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	key, err := db.GetKeyInfo(serverName, dbNum, keyName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	response := keyInfoResponse{}
	response.PageSize = defaultPageSize

	pagesCount, err := key.Values(db.NewKeyValuesQuery()).PagesCount()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	response.PagesCount = pagesCount
	response.KeyType = key.KeyType()

	return response, nil
}

//KeysSubtreeResponse is a
type KeysSubtreeResponse struct {
	Nodes []db.KeyTreeNode
	Path  []string
}

//GetKeysSubtree is a handler for getting Redis keys tree nodes
func GetKeysSubtree(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is mandatory")
	}

	var dbNum uint8
	dbNum, err := GetParamUint8("db", r)

	pathURI := r.URL.Query().Get("path")
	path := strings.Split(pathURI, "/")
	if len(pathURI) > 0 && pathURI[len(pathURI)-1] == '/' {
		path = path[:len(path)-1]
	}

	if pathURI == "" {
		path = []string{}
	}

	delimiter := ":"
	var pageSize int32
	pageSize = 100

	keyPrefix := strings.Join(path, delimiter)
	if len(path) == 0 {
		keyPrefix = "*"
	}
	node := db.KeyTreeNode{Name: keyPrefix, HasChildren: true}
	nodes, err := db.FindKeysTreeNodeChildren(serverName, dbNum, delimiter, pageSize, node)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	response := KeysSubtreeResponse{Nodes: nodes, Path: path}

	return response, nil
}

type singleValueResponse struct {
	KeyType string
	Value   db.RedisValue
}

type listValuesResponse struct {
	KeyType    string
	Values     []db.ListMember
	PageNum    int
	PagesCount int
}

type hashValuesResponse struct {
	KeyType    string
	Values     map[string]db.RedisValue
	PageNum    int
	PagesCount int
}

type setValuesResponse struct {
	KeyType    string
	Values     []db.RedisValue
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
	var dbNum uint8
	dbNum, err := GetParamUint8("db", r)

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		return nil, responds.NewBadRequestError("'key' param is required")
	}

	keyExists, err := db.KeyExists(serverName, dbNum, keyName)
	if err != nil {
		return nil, err
	}
	if !keyExists {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	pageParam := r.URL.Query().Get("page")
	pageNum := 1
	if len(pageParam) != 0 {
		if parsedPageParam, err := strconv.ParseInt(pageParam, 0, 0); err == nil {
			pageNum = int(parsedPageParam)
		}
	}

	mask := r.URL.Query().Get("mask")
	if len(mask) == 0 {
		mask = "*"
	}

	key, err := db.GetKeyInfo(serverName, dbNum, keyName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	vQuery := db.NewKeyValuesQuery()
	vQuery.PageNum = pageNum
	vQuery.Mask = mask
	vInfo := key.Values(vQuery)

	pagesCount, err := vInfo.PagesCount()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	v, err := vInfo.Values()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	switch key.KeyType() {
	case db.RedisString:
		response := singleValueResponse{}
		response.KeyType = key.KeyType()
		if rv, ok := v.(db.RedisValue); ok {
			response.Value = sanitizeBinaryValue(rv)
		}
		return response, nil

	case db.RedisList:
		response := listValuesResponse{}
		response.KeyType = key.KeyType()
		if values, ok := v.([]db.ListMember); ok {
			for i, lm := range values {
				values[i].Value = sanitizeBinaryValue(lm.Value)
			}
			response.Values = values
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil

	case db.RedisZset:
		response := zsetValuesResponse{}
		response.KeyType = key.KeyType()
		if v, ok := v.([]db.ZSetMember); ok {
			for i, zm := range v {
				v[i].Member = sanitizeBinaryValue(zm.Member)
			}
			response.Values = v
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil
	case db.RedisHash:
		response := hashValuesResponse{}
		response.KeyType = key.KeyType()
		if v, ok := v.(map[string]db.RedisValue); ok {
			for k, hv := range v {
				v[k] = sanitizeBinaryValue(hv)
			}
			response.Values = v
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil
	case db.RedisSet:
		response := setValuesResponse{}
		response.KeyType = key.KeyType()
		if v, ok := v.([]db.RedisValue); ok {
			for i, sv := range v {
				v[i] = sanitizeBinaryValue(sv)
			}
			response.Values = v
		}
		response.PageNum = pageNum
		response.PagesCount = pagesCount
		return response, nil

	default:
		return nil, fmt.Errorf("%v key has not-supported type %v", keyName, key.KeyType())
	}

}

func sanitizeBinaryValue(v db.RedisValue) db.RedisValue {
	if v.IsBinary {
		v.Value = "[binary data]"
	}

	return v
}

//DeleteKey deletes a given key
func DeleteKey(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	requestParams := server.GetURLParams(r)
	serverName := requestParams["server"]
	if len(serverName) == 0 {
		return nil, responds.NewBadRequestError("'server' param is required")
	}
	var dbNum uint8
	dbNum, err := GetParamUint8("db", r)

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		return nil, responds.NewBadRequestError("'key' param is required")
	}

	err = db.DeleteKey(serverName, dbNum, keyName)
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
	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)

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

	ex, err := db.KeyExists(serverName, dbNum, keyName)
	if err != nil {
		return nil, err
	}

	if ex {
		return nil, responds.NewConflictError(fmt.Sprintf("key %v already exists", keyName))
	}

	err = db.Set(serverName, dbNum, keyName, bodyReq.Value)
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
	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)

	decoder := json.NewDecoder(r.Body)
	var JSONReq updateStringRequest
	err = decoder.Decode(&JSONReq)
	if err != nil {
		logger.Errorf("error while parsing JSON: %v", err)
		return nil, responds.NewBadRequestError("got invalid JSON")
	}

	serverName := GetParam("server", r)
	keyName := GetParam("key", r)

	ex, err := db.KeyExists(serverName, dbNum, keyName)
	if err != nil {
		return nil, err
	}

	if !ex {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}
	if len(JSONReq.Value) == 0 {
		return nil, responds.NewBadRequestError("`Value` JSON param is missing")
	}
	err = db.Set(serverName, dbNum, keyName, JSONReq.Value)
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
	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)

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

	ex, err := db.HashKeyExists(serverName, dbNum, keyName, hashKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if ex {
		return nil, responds.NewConflictError(fmt.Sprintf("key %v already exists", keyName))
	}

	err = db.SetHashKey(serverName, dbNum, keyName, hashKey, hashValue)
	if err != nil {
		return nil, err
	}

	return "", nil
}

type updateHashValueJSONRequest struct {
	Value      string
	NewHashKey string
}

//UpdateHashValue updates an exists hash value
func UpdateHashValue(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := CheckRequiredParams([]string{"server", "key", "hashKey"}, r)
	if err != nil {
		return nil, responds.NewBadRequestError(err.Error())
	}
	serverName := GetParam("server", r)
	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)

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
	newHashKey := bodyReq.NewHashKey

	ex, err := db.HashKeyExists(serverName, dbNum, keyName, hashKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !ex {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	if len(newHashKey) > 0 && newHashKey != hashKey {
		ex, err = db.HashKeyExists(serverName, dbNum, keyName, newHashKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		if ex {
			return nil, responds.NewConflictError(fmt.Sprintf("key %s already exists for hash %s", newHashKey, keyName))
		}
	} else {
		newHashKey = hashKey
	}

	err = db.UpdateHashKey(serverName, dbNum, keyName, hashKey, newHashKey, hashValue)
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
	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)

	keyName := GetParam("key", r)
	hashKey := GetParam("hashKey", r)

	ex, err := db.HashKeyExists(serverName, dbNum, keyName, hashKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !ex {
		return nil, responds.NewNotFoundError(fmt.Sprintf("key %v doesn't exist", keyName))
	}

	err = db.DeleteHashValue(serverName, dbNum, keyName, hashKey)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	if bodyReq.Index != nil {
		err = db.InsertToListWithPos(GetParam("server", r), dbNum, GetParam("key", r), bodyReq.Value, *bodyReq.Index)
	} else {
		err = db.AppendToList(GetParam("server", r), dbNum, GetParam("key", r), bodyReq.Value)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.UpdateListValue(GetParam("server", r), dbNum, GetParam("key", r), int(index), bodyReq.Value)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.DeleteListValue(GetParam("server", r), dbNum, GetParam("key", r), int(index))
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.AddValueToSet(GetParam("server", r), dbNum, GetParam("key", r), bodyReq.Value)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.UpdateSetValue(GetParam("server", r), dbNum, GetParam("key", r), GetParam("value", r), bodyReq.Value)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.DeleteSetValue(GetParam("server", r), dbNum, GetParam("key", r), GetParam("value", r))
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.AddZSetValue(GetParam("server", r), dbNum, GetParam("key", r), bodyReq.Value, bodyReq.Score)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.UpdateZSetValueIfExists(GetParam("server", r), dbNum, GetParam("key", r), GetParam("value", r), bodyReq.Value, bodyReq.Score)
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

	var dbNum uint8
	dbNum, err = GetParamUint8("db", r)
	err = db.DeleteZSetValue(GetParam("server", r), dbNum, GetParam("key", r), GetParam("value", r))
	if err != nil {
		return nil, err
	}

	return "", nil
}

type appVersionResponse struct {
	Version string
}

//GetAppVersion returns current app version
func GetAppVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fileContents := make([]byte, 10)
	fileContents, err := ioutil.ReadFile("VERSION")
	if err != nil {
		return "", err
	}

	appVersion := string(fileContents)

	response := appVersionResponse{}
	response.Version = appVersion

	return response, nil
}
