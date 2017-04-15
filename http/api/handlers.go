package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/http/server"
	"github.com/sad0vnikov/radish/logger"
	"github.com/sad0vnikov/radish/redis/db"
)

//GetServersList is a http handler returning a list of avalable Redis instances
func GetServersList(w http.ResponseWriter, r *http.Request) {
	jsonMarshal, err := json.Marshal(config.Get().Servers)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
	}
	w.Write(jsonMarshal)
}

type getKeysByMaskResponse struct {
	Keys []string
	Page int
}

const defaultPageSize = 100

//GetKeysByMask is a http handler returning a JSON list of keys satisfying given mask
//for server with the name given in 'server' query param
func GetKeysByMask(w http.ResponseWriter, r *http.Request) {

	const pageSize = defaultPageSize

	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		respondBadRequest(w, "'server' param is mandatory")
	}
	mask := requestParams["mask"]
	if len(mask) == 0 {
		respondBadRequest(w, "'mask' param is mandatory")
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
		respondBadRequest(w, err.Error())
	}

	pageOffsetEnd := pageNumber * pageSize
	if pageOffsetEnd > len(keys) {
		pageOffsetEnd = len(keys)
	}

	pageOffsetStart := (pageNumber - 1) * pageSize
	if pageOffsetStart > len(keys) {
		respondNotFound(w)
		return
	}

	keysPage := keys[pageOffsetStart:pageOffsetEnd]

	responseContents := getKeysByMaskResponse{Keys: keysPage, Page: pageNumber}

	jsonMarshal, err := json.Marshal(responseContents)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
	}

	w.Write(jsonMarshal)
}

type keyInfoResponse struct {
	PageSize   int
	PagesCount int
	KeyType    string
}

//GetKeyInfo returns key type, values pages count and page size
func GetKeyInfo(w http.ResponseWriter, r *http.Request) {
	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		respondBadRequest(w, "'server' param is mandatory")
		return
	}

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		respondBadRequest(w, "'key' param is mandatory")
		return
	}

	key, err := db.GetKeyInfo(serverName, keyName)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
	}

	response := keyInfoResponse{}
	response.PageSize = defaultPageSize
	response.PagesCount, err = key.PagesCount(defaultPageSize)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
	}
	response.KeyType = key.KeyType()

	responseMarshal, err := json.Marshal(response)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
	}

	w.Write(responseMarshal)
}

type singleValueResponse struct {
	KeyType  string
	KeyValue string
}

type listValuesResponse struct {
	KeyType   string
	KeyValues []string
}

//GetKeyValues returns a list of key values
func GetKeyValues(w http.ResponseWriter, r *http.Request) {
	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		respondBadRequest(w, "'server' param is required")
		return
	}

	keyName := requestParams["key"]
	if len(keyName) == 0 {
		respondBadRequest(w, "'key' param is required")
		return
	}

	key, err := db.GetKeyInfo(serverName, keyName)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
		return
	}

	v, err := key.Values(1, defaultPageSize)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
		return
	}

	switch key.KeyType() {
	case db.RedisString:
		response := singleValueResponse{}
		response.KeyType = key.KeyType()
		if str, ok := v.(string); ok {
			response.KeyValue = str
			respondJSON(w, response)
		} else {
			respondInternalError(w)
		}
	default:
		respondInternalError(w)

	}

}