package api

import (
	"encoding/json"
	"net/http"

	"github.com/sad0vnikov/radish/config"
)

//GetServersList is a http handler returning a list of avalable Redis instances
func GetServersList(w http.ResponseWriter, r *http.Request) {
	jsonMarshal, err := json.Marshal(config.Get().Servers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(jsonMarshal)
}
