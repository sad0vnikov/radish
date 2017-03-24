module Model.Model exposing (Model, Server, RedisKey, initModel)

import Dict exposing (..)
import Flags exposing (Flags)

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

initModel : Flags -> Model
initModel flags =
  {
    api = { url = flags.apiUrl},
    loadedData = {
      servers = Dict.empty,
      loadedKeys = []
    },
    chosenServer = Maybe.Nothing,
    chosenKey = Maybe.Nothing
  }

type alias Server = {
  name: String,
  host: String,
  port_: Int
}

type alias RedisKey = String