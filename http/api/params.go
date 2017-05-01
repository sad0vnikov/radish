package api

import (
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"io"

	"github.com/sad0vnikov/radish/http/responds"
	"github.com/sad0vnikov/radish/http/server"
	"github.com/sad0vnikov/radish/logger"
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

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		logger.Error(err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF || err.Error() == "multipart: NextPart: EOF" {
				break
			}

			if err != nil {
				logger.Error(err)
				break
			}

			fieldName := p.FormName()
			params[fieldName] = r.FormValue(fieldName)
		}
	}

	return value
}
