package server

import (
	"log"
	"net/http"
)

//HTTPServer server configuration
type HTTPServer struct {
	Port int
}

type handler func(w http.ResponseWriter, r *http.Request)

// Init starts http server
func (server HTTPServer) Init() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//AddHandler adds a http handler
func (server HTTPServer) AddHandler(path string, h handler) {
	http.HandleFunc(path, h)
}
