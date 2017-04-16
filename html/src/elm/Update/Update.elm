module Update.Update exposing (Msg(..), update)

import Model.Model exposing (..)
import Http
import Dict exposing (..)

type Msg = NoOp 
  | ChosenServer String 
  | ServersListLoaded (Result Http.Error (Dict String Server)) 

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    ChosenServer server ->
      ({model | chosenServer = Just server}, Cmd.none)
    ServersListLoaded (Ok servers) ->
      ({model | loadedServers = (updateServersList model.loadedServers servers)}, Cmd.none)
    ServersListLoaded (Err err) ->
      ({model | error = Just ("Got error while loading servers list: " ++ (httpErrorToString err))}, Cmd.none)
    _ ->
      (model, Cmd.none)

  

updateServersList: LoadedServers -> Dict String Server -> LoadedServers
updateServersList loadedServers servers =
  {loadedServers | servers = servers}


httpErrorToString : Http.Error -> String
httpErrorToString err = 
  case err of
    Http.BadUrl string -> "Wrong request url"
    Http.Timeout -> "Request timeout"
    Http.NetworkError -> "Network error"
    Http.BadStatus _ ->  "Got error response"
    Http.BadPayload _ _ ->  "Cannot decode response"