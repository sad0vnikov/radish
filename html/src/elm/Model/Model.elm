module Model.Model exposing (Model, Server, RedisKey, LoadedServers, LoadedKeys, LoadedValues, KeyType, 
  LoadedValues(..), RedisValuesPage, RedisValue, availableKeyTypes, keyTypeName, KeyType(..), RedisValues(..), StringRedisValue, ListRedisValue, ZSetRedisValue, UserConfirmation(..), getLoadedKeyType, getChosenServerAndKey, initModel)

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
  editingValueToSave: String,
  editingScoreToSave: Int,

  isAddingValue: Bool,
  addingValue: String,
  addingHashKey: String,
  addingZSetScore: Int,

  addKeyModalShown: Bool,
  keyToAddType: String,
  keyToAddName: String
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

availableKeyTypes : List KeyType
availableKeyTypes = [StringRedisKey, HashRedisKey, SetRedisKey, ZSetRedisKey, ListRedisKey]

keyTypeName : KeyType -> String
keyTypeName keyType =
  case keyType of
    HashRedisKey -> "Hash"
    StringRedisKey -> "String Key"
    SetRedisKey -> "Set"
    ZSetRedisKey -> "ZSet"
    ListRedisKey -> "List"
    _ -> "Unsupported key type"

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
    editingValueToSave = "",
    editingScoreToSave = 0,

    isAddingValue = False,
    addingValue = "",
    addingHashKey = "",
    addingZSetScore = 0,

    addKeyModalShown = False,
    keyToAddType = keyTypeName StringRedisKey,
    keyToAddName = ""
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