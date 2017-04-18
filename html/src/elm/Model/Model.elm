module Model.Model exposing (Model, Server, RedisKey, LoadedServers, LoadedKeys, initModel)

import Dict exposing (..)
import Flags exposing (Flags)

type alias Model = {
  api: {
    url: String
  },

  loadedServers: LoadedServers,
  loadedKeys: LoadedKeys,
  keysMask: String,

  chosenServer: Maybe String,
  chosenKey: Maybe String,

  error: Maybe String
}


type alias LoadedServers = {
    servers: Dict String Server
}

type alias LoadedKeys = {
    keys: List RedisKey,
    pagesCount: Int,
    currentPage: Int
}

initModel : Flags -> Model
initModel flags =
  {
    api = { url = flags.apiUrl},
    loadedServers = {
      servers = Dict.empty
    },
    loadedKeys = {
      keys = [],
      pagesCount = 0,
      currentPage = 1
    },
    keysMask = "",
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