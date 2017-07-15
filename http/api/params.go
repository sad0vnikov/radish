package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sad0vnikov/radish/http/responds"
	"github.com/sad0vnikov/radish/http/server"
)

//CheckRequiredParams receives a list of required URL and HTTP params and returns error if any param is missing
func CheckRequiredParams(params []string, r *http.Request) error {

	requestParams := server.GetURLParams(r)
	for _, p := range params {
		pURLValue := requestParams[p]
		pQueryValue := r.URL.Query().Get(p)

		if len(pURLValue) == 0 && len(pQueryValue) == 0 {
			return responds.NewBadRequestError(fmt.Sprintf("'%v' param is required", p))
		}
	}
	return nil
}

//GetParam returns a URL or Query param value
func GetParam(param string, r *http.Request) string {
	urlParams := server.GetURLParams(r)
	value := urlParams[param]
	if len(value) == 0 {
		value = r.URL.Query().Get(param)
	}

	return value
}

//GetParamUint8 returns a URL or Query param value as uint8
func GetParamUint8(param string, r *http.Request) (uint8, error) {
	urlParams := server.GetURLParams(r)
	value := urlParams[param]
	if len(value) == 0 {
		value = r.URL.Query().Get(param)
	}

	conv, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}

	return uint8(conv), nil
}
