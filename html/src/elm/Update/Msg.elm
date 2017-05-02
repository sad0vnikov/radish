module Update.Msg exposing (Msg(..))

import Http
import Dict exposing (Dict)
import Model.Model exposing (LoadedKeys, Server, RedisKey, RedisValue, StringRedisValue, LoadedValues)

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
  | KeyDeleted (Result Http.Error ())
  | ValueDeletionConfirm String
  | ValueDeletionConfirmed String
  | ValueDeleted (Result Http.Error ())
  | UserConfirmation
  | UserConfirmationCancel
  | ValueToEditSelected (String, StringRedisValue)
  | ZSetValueToEditSelected (String, StringRedisValue, Int)  
  | ValueEditingCanceled
  | EditedValueChanged String
  | EditedScoreChanged String
  | ValueUpdateInitialized StringRedisValue
  | ValueUpdated (Result Http.Error ())
  | AddingValueStart
  | AddingValueCancel
  | AddingValueChanged String
  | AddingHashKeyChanged String
  | AddingZSetScoreChanged String
  | AddingValueInitialized
  | ValueAdded (Result Http.Error ())