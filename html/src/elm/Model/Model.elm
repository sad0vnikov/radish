module Model.Model exposing (Model, Server, RedisKey, LoadedServers, LoadedKeys, LoadedValues, KeyType, 
  LoadedValues(..), RedisValuesPage, RedisValue, KeyType(..), RedisValues(..), StringRedisValue, ListRedisValue, ZSetRedisValue, UserConfirmation(..), getLoadedKeyType, getChosenServerAndKey, initModel)

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

  waitingForConfirmation: Maybe UserConfirmation,
  editingValue: Maybe (RedisKey, String),
  editingValueToSave: String
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
  | ListRedisValues (List ListRedisValue )
  | HashRedisValues (Dict String StringRedisValue)
  | SetRedisValues (List StringRedisValue)
  | ZSetRedisValues (List ZSetRedisValue)

type alias StringRedisValue = String

type alias ZSetRedisValue = {
    score: Int,
    value: String
}

type alias ListRedisValue = {
    index: Int,
    value: String
}

type UserConfirmation = KeyDeletion String | ValueDeletion String


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
    waitingForConfirmation = Nothing,
    editingValue = Nothing,
    editingValueToSave = ""
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