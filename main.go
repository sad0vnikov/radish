package main

import (
	"./config"
	"./http/api"
	"./http/server"
)

func main() {
	_, err := config.Load()
	if err != nil {
		panic(err)
	}

	server := server.HTTPServer{Port: 8080}
	server.Init()
	server.AddHandler(api.Version()+"/servers", api.GetServersList)

}
