package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sad0vnikov/radish/http/responds"
	"github.com/sad0vnikov/radish/logger"
)

//HTTPServer server configuration
type HTTPServer struct {
	Port int
}

type apiHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

var router = mux.NewRouter()

// Init starts http server
func (server HTTPServer) Init() {
	loggingRouter := handlers.LoggingHandler(logger.GetOutput(), router)
	err := http.ListenAndServe(":8080", loggingRouter)
	if err != nil {
		logger.Critical(err)
	}
}

//ServeStatic turns on serving Radish panel static files
func (server HTTPServer) ServeStatic() {
	fs := http.FileServer(http.Dir("html/dist"))
	router.PathPrefix("/").Handler(fs)
}

//AddHandler adds a http handler
func (server HTTPServer) AddHandler(method, path string, h apiHandler) {
	router.HandleFunc(
		"/api/"+path,
		func(w http.ResponseWriter, r *http.Request) {
			resp, err := h(w, r)
			if _, ok := err.(*responds.APINotFoundError); ok {
				responds.RespondNotFound(w)
				return
			}
			if brerr, ok := err.(*responds.APIBadRequestError); ok {
				responds.RespondBadRequest(w, brerr.Error())
				return
			}
			if err != nil {
				responds.RespondInternalError(w)
				return
			}

			responds.RespondJSON(w, resp)
		}).
		Methods(method)
}

//GetURLParams returns request params from given HTTP request
func GetURLParams(request *http.Request) map[string]string {
	return mux.Vars(request)
}
