package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sad0vnikov/radish/logger"
)

func respondInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Sorry! Something went wrong")
}

func respondBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}

func respondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func respondJSON(w http.ResponseWriter, response interface{}) {
	responseMarshal, err := json.Marshal(response)
	if err != nil {
		logger.Error(err)
		respondInternalError(w)
		return
	}

	w.Write(responseMarshal)
}
