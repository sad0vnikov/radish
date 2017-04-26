module Model.Model exposing (Model, Server, RedisKey, LoadedServers, LoadedKeys, LoadedValues, KeyType, 
  LoadedValues(..), RedisValuesPage, RedisValue, KeyType(..), RedisValues(..), StringRedisValue, ZSetRedisValue, UserConfirmation(..), getLoadedKeyType, getChosenServerAndKey, initModel)

import Dict exposing (..)
import Flags exposing (Flags)

type alias Model = {
  api: {
    url: String
  },

  loadedServers: LoadedServers,
  loadedKeys: LoadedKeys,
  loadedValues: LoadedValues,
  
  keysMask: String,

  chosenServer: Maybe String,
  chosenKey: Maybe String,

  error: Maybe String,

  waitingForConfirmation: Maybe UserConfirmation
}


type alias LoadedServers = {
    servers: Dict String Server
}

type alias LoadedKeys = {
    keys: List RedisKey,
    pagesCount: Int,
    currentPage: Int
}

type LoadedValues = MultipleRedisValues RedisValuesPage | SingleRedisValue RedisValue

type alias RedisValuesPage =  {
    values: RedisValues,
    pagesCount: Int,
    currentPage: Int,
    keyType: KeyType
}

type alias RedisValue = {
    value: StringRedisValue,
    keyType: KeyType
}

type KeyType = StringRedisKey | HashRedisKey | SetRedisKey | ZSetRedisKey | ListRedisKey | UnknownRedisKey

type RedisValues = StringRedisValue String 
  | ListRedisValues (List StringRedisValue )
  | HashRedisValues (Dict String StringRedisValue)
  | SetRedisValues (List StringRedisValue)
  | ZSetRedisValues (List ZSetRedisValue)

type alias StringRedisValue = String

type alias ZSetRedisValue = {
    score: Int,
    value: String
}

type UserConfirmation = KeyDeletion String


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
    loadedValues = SingleRedisValue <| RedisValue "" StringRedisKey,
    keysMask = "",
    chosenServer = Maybe.Nothing,
    chosenKey = Maybe.Nothing,
    error = Maybe.Nothing,
    waitingForConfirmation = Nothing
  }

getLoadedKeyType : LoadedValues -> KeyType
getLoadedKeyType loadedValues =
  case loadedValues of
    MultipleRedisValues values -> values.keyType
    SingleRedisValue value -> value.keyType

type alias Server = {
  name: String,
  host: String,
  port_: Int
}

type alias RedisKey = String

getChosenServerAndKey : Model -> Maybe (String, String)
getChosenServerAndKey model =
    Maybe.map2 (\key server -> (key,server)) model.chosenServer model.chosenKey