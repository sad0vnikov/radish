package api

import (
	"encoding/json"
	"net/http"

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

//GetKeysByMask is a http handler returning a JSON list of keys satisfying given mask
//for server with the name given in 'server' query param
func GetKeysByMask(w http.ResponseWriter, r *http.Request) {

	requestParams := server.GetURLParams(r)

	serverName := requestParams["server"]
	if len(serverName) == 0 {
		respondBadRequest(w, "'server' param is mandatory")
	}
	mask := requestParams["mask"]
	if len(mask) == 0 {
		respondBadRequest(w, "'mask' param is mandatory")
	}

	keys, err := db.FindKeysByMask(serverName, mask)
	if err != nil {
		respondBadRequest(w, err.Error())
	}

	jsonMarshal, err := json.Marshal(keys)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
	}

	w.Write(jsonMarshal)
}
