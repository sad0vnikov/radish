module Model.Model exposing (Model, Server, RedisKey, LoadedServers, LoadedKeys, LoadedValues, KeyType, KeysTreeNode(..), LoadedKeysSubtree, UnfoldKeysTreeNodeInfo, CollapsedKeysTreeNodeInfo, KeysTreeLeafInfo,
  LoadedValues(..), RedisValuesPage, RedisValue, emptyKeysSubtree, availableKeyTypes, keyTypeName, keyTypeAlias, keyTypeFromAlias, KeysViewType(..), KeyType(..), RedisValues(..), StringRedisValue, 
  ListRedisValue, ZSetRedisValue, SetRedisValue, HashRedisValue, UserConfirmation(..), getChosenServer, getLoadedKeyType, getChosenServerAndKey, initModel)

import Dict exposing (..)
import Flags exposing (Flags)

type alias Model = {
  appInfo: {
    version: String
  },

  api: {
    url: String
  },

  windowSize: {
    width: Int,
    height: Int
  },

  loadedServers: LoadedServers,
  loadedKeys: LoadedKeys,
  loadedKeysTree: LoadedKeysSubtree,
  loadedValues: LoadedValues,
  
  keysMask: String,

  chosenServer: Maybe String,
  chosenDatabaseNum: Int,  
  chosenKeysViewType: KeysViewType,
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
  aboutWindowShown: Bool,  
  keyToAddType: String,
  keyToAddName: String
}


type alias LoadedServers = {
    servers: Dict String Server
}

type KeysViewType = KeysListView | KeysTreeView
type alias LoadedKeys = {
    keys: List RedisKey,
    pagesCount: Int,
    currentPage: Int
}

type alias LoadedKeysSubtree = {
  path: List String,
  loadedNodes: List KeysTreeNode
}

emptyKeysSubtree : List String -> LoadedKeysSubtree
emptyKeysSubtree path =
  LoadedKeysSubtree path []

type KeysTreeNode = UnfoldKeyTreeNode UnfoldKeysTreeNodeInfo | CollapsedKeyTreeNode CollapsedKeysTreeNodeInfo | KeysTreeLeaf KeysTreeLeafInfo

type alias UnfoldKeysTreeNodeInfo = {
    path: List String,
    name: String,
    loadedChildren: LoadedKeysSubtree
}

type alias CollapsedKeysTreeNodeInfo = {
    path : List String,
    name: String
}

type alias KeysTreeLeafInfo = {
    key: RedisKey,
    name: String
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
    keyType: KeyType,
    isBinary: Bool
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

keyTypeAlias : KeyType -> String
keyTypeAlias keyType =
  case keyType of
    HashRedisKey -> "hash"
    StringRedisKey -> "string"
    SetRedisKey -> "set"
    ZSetRedisKey -> "zset"
    ListRedisKey -> "list"
    _ -> "unknown"

keyTypeFromAlias : String -> KeyType
keyTypeFromAlias slug =
  case slug of
    "hash" -> HashRedisKey
    "string" -> StringRedisKey
    "set" -> SetRedisKey
    "zset" -> ZSetRedisKey
    "list" -> ListRedisKey
    _ -> UnknownRedisKey

type RedisValues = StringRedisValue RedisValue
  | ListRedisValues (List ListRedisValue )
  | HashRedisValues (Dict String HashRedisValue)
  | SetRedisValues (List SetRedisValue)
  | ZSetRedisValues (List ZSetRedisValue)

type alias StringRedisValue = String

type alias ZSetRedisValue = {
    score: Int,
    value: String,
    isBinary: Bool
}

type alias SetRedisValue = {
    value: String,
    isBinary: Bool
}

type alias HashRedisValue = {
    value: String,
    isBinary: Bool
}

type alias ListRedisValue = {
    index: Int,
    value: String,
    isBinary: Bool
}

type UserConfirmation = KeyDeletion String | ValueDeletion String


initModel : Flags -> Model
initModel flags =
  {
    appInfo = { version = flags.appVersion },
    api = { url = flags.apiUrl},
    loadedServers = {
      servers = Dict.empty
    },
    chosenKeysViewType = KeysListView,
    loadedKeys = {
      keys = [],
      pagesCount = 0,
      currentPage = 1
    },
    windowSize = {width = 0, height = 0},
    loadedKeysTree = emptyKeysSubtree [],
    loadedValues = SingleRedisValue <| RedisValue "" StringRedisKey False,
    keysMask = "",
    chosenServer = Maybe.Nothing,   
    chosenDatabaseNum = 0, 
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
    aboutWindowShown = False,    
    keyToAddType = keyTypeAlias StringRedisKey,
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
  port_: Int,
  databasesCount: Int,
  connectionCheckPassed: Bool
}

type alias RedisKey = String

getChosenServerAndKey : Model -> Maybe (String, String)
getChosenServerAndKey model =
    Maybe.map2 (\key server -> (key,server)) model.chosenServer model.chosenKey

getChosenServer : Model -> Maybe Server
getChosenServer model =
  case model.chosenServer of
    Just server -> 
      Dict.get server model.loadedServers.servers
    Nothing -> Nothing