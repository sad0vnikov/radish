module Update.Update exposing (update)

import Model.Model exposing (..)
import Http
import Dict exposing (..)
import View.Toastr as Toastr
import Update.Msg exposing (Msg(..))
import Command.Values exposing (..)
import Command.Keys exposing (..)
import View.ConfirmationDialog exposing (..)

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
        updatedModel = {model | keysMask = mask, loadedKeys = updateKeysPage model.loadedKeys 1}
      in
        (updatedModel, getKeysPage updatedModel)
    KeysPageChanged pageNum ->
      let 
        updatedModel = {model | loadedKeys = updateKeysPage model.loadedKeys pageNum}
      in
        (updatedModel, getKeysPage updatedModel)
    KeyChosen key ->
      let
        updatedModel = {model | chosenKey = Just key}
      in
        (updatedModel, getKeyValues updatedModel)
    UserConfirmation ->
      case model.waitingForConfirmation of
        Just value ->
          getConfirmedMessage model value
        Nothing ->
          (model, Cmd.none)
    KeyValuesLoaded (Ok values) ->
      ({model | loadedValues = values}, Cmd.none)
    KeyValuesLoaded (Err err) ->
      let
        errorStr = "Got error while loading keys list: " ++ (httpErrorToString err)
      in
        (model, Toastr.toastError errorStr)
    KeyDeletionConfirm key ->
      case model.chosenKey of
        Just chosenKey ->
          ({model | waitingForConfirmation = Just (KeyDeletion chosenKey)}, showConfirmationDialog "Do you really want to delete this key?")
        Nothing ->
          (model, Cmd.none)
    KeyDeletionConfirmed key ->
      (model, deleteKey key model) 
    KeyDeleted (Ok response) ->
      ({model | chosenKey = Nothing}, getKeysPage model)
    KeyDeleted (Err err) ->
      (model, Toastr.toastError <| "Got error while deleting key: " ++ (httpErrorToString err) )
    ValueDeletionConfirm value ->
      ({model | waitingForConfirmation = Just (ValueDeletion value)}, showConfirmationDialog "Do you really want to delete this value?")
    ValueDeletionConfirmed value ->
      (model, deleteValue model value)
    ValueDeleted (Ok response) ->
      (model, getKeyValues model) 
    ValueDeleted (Err err) ->
      (model, Toastr.toastError <| "Got error while deleting value: " ++ (httpErrorToString err))
    ValueToEditSelected value ->
      case model.chosenKey of
        Just key -> ({model | editingValue = Just (key, value), editingValueToSave = value}, Cmd.none)
        Nothing -> (model, Cmd.none)
    EditedValueChanged value ->
      ({model | editingValueToSave = value}, Cmd.none)
    ValueEditingCanceled ->
      ({model | editingValue = Nothing}, Cmd.none)
    ValueUpdateInitialized value ->
      case model.chosenKey of
        Just key -> (model, updateValue model value model.editingValueToSave)
        Nothing -> (model, Cmd.none)
    ValueUpdated (Ok response) ->
      ({model | editingValue = Nothing}, getKeyValues model)
    ValueUpdated (Err err) ->
      (model, Toastr.toastError <| "Error while updating value: " ++ (httpErrorToString err))
    _ ->
      (model, Cmd.none)

getConfirmedMessage : Model -> UserConfirmation -> (Model, Cmd Msg)
getConfirmedMessage model confirmation =
  case confirmation of
    KeyDeletion key ->
      update (KeyDeletionConfirmed key) model
    ValueDeletion value ->
      update (ValueDeletionConfirmed value) model


updateServersList: LoadedServers -> Dict String Server -> LoadedServers
updateServersList loadedServers servers =
  {loadedServers | servers = servers}

updateKeysPage : LoadedKeys -> Int -> LoadedKeys
updateKeysPage loadedKeys newPage =
  {loadedKeys | currentPage = newPage }


httpErrorToString : Http.Error -> String
httpErrorToString err = 
  case err of
    Http.BadUrl string -> "Wrong request url"
    Http.Timeout -> "Request timeout"
    Http.NetworkError -> "Network error"
    Http.BadStatus _ ->  "Got error response"
    Http.BadPayload _ _ ->  "Cannot decode response"