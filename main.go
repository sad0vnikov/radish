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
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys", api.GetKeysByMask)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/{key}/info", api.GetKeyInfo)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/{key}/values", api.GetKeyValues)
	server.AddHandler("DELETE", api.Version()+"/servers/{server}/keys/{key}", api.DeleteKey)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/strings/{key}/values", api.AddStringValue)
	server.AddHandler("PUT", api.Version()+"/servers/{server}/keys/strings/{key}/values", api.UpdateStringValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/lists/{key}/values", api.AddListValue)
	server.AddHandler("PUT", api.Version()+"/servers/{server}/keys/lists/{key}/values", api.UpdateListValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/hashes/{key}/values", api.AddHashValue)
	server.AddHandler("PUT", api.Version()+"/servers/{server}/keys/hashes/{key}/values", api.UpdateHashValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/sets/{key}/values", api.AddSetValue)
	server.AddHandler("PUT", api.Version()+"/servers/{server}/keys/sets/{key}/values", api.UpdateSetValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/zsets/{key}/values", api.AddZSetValue)
	server.AddHandler("PUT", api.Version()+"/servers/{server}/keys/zsets/{key}/values", api.UpdateZSetValue)

	server.ServeStatic()
	server.Init()

}
