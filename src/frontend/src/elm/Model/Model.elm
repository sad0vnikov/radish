module Model.Model exposing (Model, Server, RedisKey)

import Dict exposing (..)

type alias Model = {
  api: {
    url: String
  },
  loadedData: {
    servers: Dict String Server,
    loadedKeys: List RedisKey
  }, 
  chosenServer: Maybe String,
  chosenKey: Maybe String
}



type alias Server = {
  name: String,
  host: String,
  port_: Int
}

type alias RedisKey = String