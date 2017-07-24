module Update.Msg exposing (Msg(..))

import Http
import Window
import Dict exposing (Dict)
import Model.Model exposing (LoadedKeys, Server, RedisKey, RedisValue, StringRedisValue, LoadedValues, KeyType, KeysTreeNode(..), LoadedKeysSubtree, CollapsedKeysTreeNodeInfo, UnfoldKeysTreeNodeInfo)

type Msg = NoOp 
  | ChosenServer String 
  | DatabaseChosen Int
  | KeysListUpdateInitialized
  | ServersListLoaded (Result Http.Error (Dict String Server))
  | KeysPageLoaded (Result Http.Error LoadedKeys)
  | KeysMaskChanged String
  | ValuesMaskChanged String
  | ShowValuesFilter
  | HideValuesFilter
  | KeysPageChanged Int
  | KeyChosen RedisKey
  | ValuesListUpdateInitialized
  | ValuesPageChanged Int
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
  | HashValueToEditSelected (String, StringRedisValue)
  | ZSetValueToEditSelected (String, StringRedisValue, Int)  
  | ValueEditingCanceled
  | EditedValueChanged String
  | EditedHashKeyChanged String
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
  | ShowAddKeyModal
  | CloseAddKeyModal
  | KeyToAddTypeChanged String
  | KeyToAddNameChanged String
  | AddNewKey
  | ListKeyViewChosen
  | TreeKeyViewChosen
  | KeysTreeSubtreeLoaded (Result Http.Error LoadedKeysSubtree)
  | KeysTreeCollapsedNodeClick CollapsedKeysTreeNodeInfo
  | KeysTreeUnfoldNodeClick UnfoldKeysTreeNodeInfo
  | LogoClick
  | AboutWindowOpen
  | AboutWindowClose
  | WindowResized (Int, Int)