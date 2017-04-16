package responds

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sad0vnikov/radish/logger"
)

//APINotFoundError is a NotFound error returned by Api endpoint handler
type APINotFoundError struct {
	msg string
}

func (err APINotFoundError) Error() string {
	return err.msg
}

//NewNotFoundError returns a new API not found error
func NewNotFoundError(msg string) error {
	return &APINotFoundError{msg: msg}
}

//APIBadRequestError is a Bad Request error returned by Api endpoint handler
type APIBadRequestError struct {
	msg string
}

func (err APIBadRequestError) Error() string {
	return err.msg
}

//NewBadRequestError returns a new BadRequest error
func NewBadRequestError(msg string) error {
	return &APIBadRequestError{msg}
}

//RespondInternalError responds with 500 Internal Error HTTP status
func RespondInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

//RespondBadRequest responds with 400 Bad Request HTTP status
func RespondBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}

//RespondNotFound responds with 404 Not Found HTTP status
func RespondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

//RespondJSON writes JSON to http output
func RespondJSON(w http.ResponseWriter, response interface{}) {
	responseMarshal, err := json.Marshal(response)
	if err != nil {
		logger.Error(err)
		RespondInternalError(w)
		return
	}

	w.Write(responseMarshal)
}
