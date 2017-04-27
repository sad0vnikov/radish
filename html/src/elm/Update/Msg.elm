module Update.Msg exposing (Msg(..))

import Http
import Dict exposing (Dict)
import Model.Model exposing (LoadedKeys, Server, RedisKey, RedisValue, LoadedValues)

type Msg = NoOp 
  | ChosenServer String 
  | ServersListLoaded (Result Http.Error (Dict String Server)) 
  | KeysPageLoaded (Result Http.Error LoadedKeys)
  | KeysMaskChanged String
  | KeysPageChanged Int
  | KeyChosen RedisKey
  | KeyValuesLoaded (Result Http.Error LoadedValues)
  | KeyDeletionConfirm RedisKey
  | KeyDeletionConfirmed RedisKey
  | KeyDeleted (Result Http.Error String)
  | ValueDeletionConfirm String
  | ValueDeletionConfirmed String
  | ValueDeleted (Result Http.Error String)
  | UserConfirmation
  | UserConfirmationCancel