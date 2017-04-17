module Update.Update exposing (update)

import Model.Model exposing (..)
import Http
import Dict exposing (..)
import View.Toastr as Toastr
import Update.Msg exposing (Msg(..))
import Command.Keys exposing (..)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    ChosenServer server ->
      let
        updatedModel = {model | chosenServer = Just server}
      in
        (updatedModel, getKeysPage updatedModel)

    ServersListLoaded (Ok servers) ->
      ({model | loadedServers = (updateServersList model.loadedServers servers)}, Cmd.none)
    ServersListLoaded (Err err) ->
      let
        errorStr = "Got error while loading servers list: " ++ (httpErrorToString err)
      in
        (model, Toastr.toastError errorStr)
    KeysPageLoaded (Ok keys) ->
      ({model | loadedKeys = keys}, Cmd.none)
    KeysPageLoaded (Err err) ->
      let
        errorStr = "Got error while loading keys list: " ++ (httpErrorToString err)
      in
        (model, Toastr.toastError errorStr)
    KeysMaskChanged mask ->
      let 
        updatedModel = {model | keysMask = mask}
      in
        (updatedModel, getKeysPage updatedModel)
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