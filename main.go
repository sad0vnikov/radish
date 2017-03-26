package main

import (
	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/http/api"
	"github.com/sad0vnikov/radish/http/server"
	"github.com/sad0vnikov/radish/logger"
)

func main() {
	logger.Info("init app...")
	var configLoader config.Loader

	configPath := "config.json"
	configLoader = config.JSONFileConfigLoader{Path: configPath}
	logger.Infof("read config from %v", configPath)

	_, err := configLoader.Load()

	if err != nil {
		panic(err)
	}

	server := server.HTTPServer{Port: 8080}
	server.AddHandler("GET", api.Version()+"/servers", api.GetServersList)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/bymask/{mask}", api.GetKeysByMask)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/bymask/{mask}/page/{page}", api.GetKeysByMask)

	server.Init()

}
