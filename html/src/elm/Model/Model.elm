module Model.Model exposing (Model, Server, RedisKey, LoadedServers, initModel)

import Dict exposing (..)
import Flags exposing (Flags)

type alias Model = {
  api: {
    url: String
  },

  loadedServers: LoadedServers,
  loadedKeys: List RedisKey,

  chosenServer: Maybe String,
  chosenKey: Maybe String,

  error: Maybe String
}


type alias LoadedServers = {
    servers: Dict String Server
}

initModel : Flags -> Model
initModel flags =
  {
    api = { url = flags.apiUrl},
    loadedServers = {
      servers = Dict.empty
    },
    loadedKeys = [],
    chosenServer = Maybe.Nothing,
    chosenKey = Maybe.Nothing,
    error = Maybe.Nothing
  }

type alias Server = {
  name: String,
  host: String,
  port_: Int
}

type alias RedisKey = String