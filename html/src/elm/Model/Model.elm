module Model.Model exposing (Model, Server, RedisKey, LoadedServers, LoadedKeys, LoadedValues, KeyType, KeysTreeNode(..), LoadedKeysSubtree, UnfoldKeysTreeNodeInfo, CollapsedKeysTreeNodeInfo,
  LoadedValues(..), RedisValuesPage, RedisValue, availableKeyTypes, keyTypeName, keyTypeAlias, keyTypeFromAlias, KeysViewType(..), KeyType(..), RedisValues(..), StringRedisValue, ListRedisValue, ZSetRedisValue, UserConfirmation(..), getLoadedKeyType, getChosenServerAndKey, initModel)

import Dict exposing (..)
import Flags exposing (Flags)

type alias Model = {
  api: {
    url: String
  },

  loadedServers: LoadedServers,
  loadedKeys: LoadedKeys,
  loadedKeysTree: LoadedKeysSubtree,
  loadedValues: LoadedValues,
  
  keysMask: String,

  chosenServer: Maybe String,
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
  nodesCount: Int,
  path: List String,
  loadedNodes: List KeysTreeNode
}

type KeysTreeNode = UnfoldKeyTreeNode UnfoldKeysTreeNodeInfo | CollapsedKeyTreeNode CollapsedKeysTreeNodeInfo | KeysTreeLeaf RedisKey

type alias UnfoldKeysTreeNodeInfo = {
    name: String,
    childrenCount: Int,
    loadedChildren: List LoadedKeysSubtree
}

type alias CollapsedKeysTreeNodeInfo = {
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
    chosenKeysViewType = KeysListView,
    loadedKeys = {
      keys = [],
      pagesCount = 0,
      currentPage = 1
    },
    loadedKeysTree = {
      nodesCount = 0,
      path = [],
      loadedNodes = []
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
  port_: Int
}

type alias RedisKey = String

getChosenServerAndKey : Model -> Maybe (String, String)
getChosenServerAndKey model =
    Maybe.map2 (\key server -> (key,server)) model.chosenServer model.chosenKey