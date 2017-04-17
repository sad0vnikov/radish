module Update.Msg exposing (Msg(..))

import Http
import Dict exposing (Dict)
import Model.Model exposing (RedisKey, Server)

type Msg = NoOp 
  | ChosenServer String 
  | ServersListLoaded (Result Http.Error (Dict String Server)) 
  | KeysPageLoaded (Result Http.Error (List RedisKey))
  | KeysMaskChanged String