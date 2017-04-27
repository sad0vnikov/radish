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
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/{key}/delete", api.DeleteKey)

	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/strings/{key}/values/add", api.AddStringValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/strings/{key}/values/update", api.UpdateStringValue)

	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/lists/{key}/values/add", api.AddListValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/lists/{key}/values/update", api.UpdateListValue)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/lists/{key}/values/{index}/delete", api.DeleteListValue)

	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/hashes/{key}/values/add", api.AddHashValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/hashes/{key}/values/update", api.UpdateHashValue)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/hashes/{key}/values/{hashKey}/delete", api.DeleteHashValue)

	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/sets/{key}/values/add", api.AddSetValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/sets/{key}/values/update", api.UpdateSetValue)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/sets/{key}/values/{value}/delete", api.DeleteSetValue)

	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/zsets/{key}/values/add", api.AddZSetValue)
	server.AddHandler("POST", api.Version()+"/servers/{server}/keys/zsets/{key}/values/update", api.UpdateZSetValue)
	server.AddHandler("GET", api.Version()+"/servers/{server}/keys/zsets/{key}/values/{value}/delete", api.DeleteZSetValue)

	server.ServeStatic()
	server.Init()

}
