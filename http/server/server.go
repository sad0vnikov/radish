package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sad0vnikov/radish/logger"
)

//HTTPServer server configuration
type HTTPServer struct {
	Port int
}

type handler func(w http.ResponseWriter, r *http.Request)

var router = mux.NewRouter()

// Init starts http server
func (server HTTPServer) Init() {
	loggingRouter := handlers.LoggingHandler(logger.GetOutput(), router)
	err := http.ListenAndServe(":8080", loggingRouter)
	if err != nil {
		logger.Critical(err)
	}
}

//AddHandler adds a http handler
func (server HTTPServer) AddHandler(method, path string, h handler) {
	router.HandleFunc("/api/"+path, h).
		Methods(method)
}
