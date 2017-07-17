module Update.Update exposing (update)

import Model.Model exposing (..)
import Http
import String exposing (toInt)
import Dict exposing (..)
import View.Toastr as Toastr
import Update.Msg exposing (Msg(..))
import Command.Values exposing (..)
import Command.Keys exposing (..)
import View.ConfirmationDialog exposing (..)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    WindowResized (width, height) ->
      ({model | windowSize = {width = width, height = height}}, Cmd.none)
    ChosenServer server ->
      let
        updatedModel = {model | chosenServer = Just server}
      in
        (updatedModel, getKeysPage updatedModel)
    DatabaseChosen dbNum ->
      let
        updatedModel = {model | chosenDatabaseNum = dbNum, chosenKey = Nothing}
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
        updatedModel = {model | chosenKey = Just key, editingValue = Nothing, isAddingValue = False}
      in
        (updatedModel, getKeyValues updatedModel 1)
    ValuesPageChanged page ->
      (model, getKeyValues model page)
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
      (model, getKeyValues model 1) 
    ValueDeleted (Err err) ->
      (model, Toastr.toastError <| "Got error while deleting value: " ++ (httpErrorToString err))
    ValueToEditSelected (valueReference, currentValue) ->
      case model.chosenKey of
        Just key -> ({model | editingValue = Just (key, valueReference), editingValueToSave = currentValue}, Cmd.none)
        Nothing -> (model, Cmd.none)
    ZSetValueToEditSelected (valueReference, currentValue, currentScore) ->
      let 
        modelWithScore = {model | editingScoreToSave = currentScore}
      in
        update (ValueToEditSelected (valueReference, currentValue)) modelWithScore

    EditedValueChanged value ->
      ({model | editingValueToSave = value}, Cmd.none)
    EditedScoreChanged score ->
      let 
        convertedScore = String.toInt score
      in
        case convertedScore of
          Ok score -> ({model | editingScoreToSave = score}, Cmd.none)
          Err _ -> (model, Toastr.toastError "Score should be number")
       
    ValueEditingCanceled ->
      ({model | editingValue = Nothing}, Cmd.none)
    ValueUpdateInitialized value ->
      case model.chosenKey of
        Just key -> (model, updateValue model value model.editingValueToSave)
        Nothing -> (model, Cmd.none)
    ValueUpdated (Ok response) ->
      ({model | editingValue = Nothing}, getKeyValues model 1)
    ValueUpdated (Err err) ->
      (model, Toastr.toastError <| "Error while updating value: " ++ (httpErrorToString err))

    AddingValueStart ->
      ({model | isAddingValue = True, addingValue = "", addingHashKey = ""}, Cmd.none)
    AddingValueCancel ->
      ({model | isAddingValue = False, addingValue = "", addingHashKey = ""}, Cmd.none)
    AddingValueChanged value ->
      ({model | addingValue = value}, Cmd.none)
    AddingHashKeyChanged value ->
      ({model | addingHashKey = value}, Cmd.none)
    AddingZSetScoreChanged stringValue ->
      case String.toInt stringValue of
        Ok value -> ({model | addingZSetScore = value}, Cmd.none)
        Err _ -> (model, Cmd.none)
    AddingValueInitialized ->
      (model, addValueForChosenKey model)
    ValueAdded (Ok response) ->
      ({model | isAddingValue = False, addKeyModalShown = False}, getKeyValues model 1)
    ValueAdded (Err err) ->
      (model, Toastr.toastError <| "Error while adding value: " ++ (valueAddingErrorToString err))
    ShowAddKeyModal ->
      ({model | addKeyModalShown = True}, Cmd.none)
    CloseAddKeyModal ->
      ({model | addKeyModalShown = False}, Cmd.none)
    KeyToAddTypeChanged keyType ->
      ({model | keyToAddType = keyType}, Cmd.none)      
    KeyToAddNameChanged keyName ->
      ({model | keyToAddName = keyName}, Cmd.none)   
    AddNewKey ->
      ({model | chosenKey = Just model.keyToAddName}, addKey model)
    ListKeyViewChosen ->
      let
        updatedModel = {model | chosenKeysViewType = KeysListView}
      in
        (updatedModel, getKeysPage updatedModel)
    TreeKeyViewChosen ->
      let
        updatedModel = {model | chosenKeysViewType = KeysTreeView}
      in
        (updatedModel, getKeysSubtree updatedModel [])
    KeysTreeSubtreeLoaded (Err err) ->
      let
        errorStr = "Got error while loading keys subtree: " ++ (httpErrorToString err)
      in
        (model, Toastr.toastError errorStr)
    KeysTreeSubtreeLoaded (Ok loadedSubtree) ->
      ({model | loadedKeysTree = updateKeysTree loadedSubtree model.loadedKeysTree}, Cmd.none)
    KeysTreeCollapsedNodeClick node ->
      let 
        subtreeToLoadPath = node.path ++ (List.singleton node.name)
      in
        (model, getKeysSubtree model subtreeToLoadPath)
    KeysTreeUnfoldNodeClick node ->
      ({model | loadedKeysTree = collapseKeysTreeNode model.loadedKeysTree node}, Cmd.none)
    LogoClick -> 
      update AboutWindowOpen model
    AboutWindowOpen ->
      ({model | aboutWindowShown = True}, Cmd.none)
    AboutWindowClose->
      ({model | aboutWindowShown = False}, Cmd.none)
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


updateKeysTree : LoadedKeysSubtree -> LoadedKeysSubtree -> LoadedKeysSubtree
updateKeysTree loadedSubtree currentSubtree =
  if List.isEmpty loadedSubtree.path then
    {currentSubtree | loadedNodes = (currentSubtree.loadedNodes ++ loadedSubtree.loadedNodes)}
  else
    {currentSubtree | loadedNodes = List.map (updateSubtreeLoadedNode loadedSubtree) currentSubtree.loadedNodes}

updateSubtreeLoadedNode : LoadedKeysSubtree -> KeysTreeNode -> KeysTreeNode
updateSubtreeLoadedNode loadedSubtree node =
  case node of
    UnfoldKeyTreeNode unfoldKeyInfo ->
      if Just unfoldKeyInfo.name == List.head loadedSubtree.path then
        let 
          subtreeWithCroppedPath = {loadedSubtree | path = List.drop 1 loadedSubtree.path}
        in
          UnfoldKeyTreeNode {unfoldKeyInfo | loadedChildren = updateKeysTree subtreeWithCroppedPath unfoldKeyInfo.loadedChildren}
      else
        node  
    CollapsedKeyTreeNode keyInfo ->
      if List.length loadedSubtree.path == 1 && Just keyInfo.name == List.head loadedSubtree.path then
        let 
           subtreeWithCroppedPath = {loadedSubtree | path = []}
        in
          UnfoldKeyTreeNode <| UnfoldKeysTreeNodeInfo keyInfo.path keyInfo.name (updateKeysTree subtreeWithCroppedPath (emptyKeysSubtree []))
      else
        node
    _ ->
      node

collapseKeysTreeNode : LoadedKeysSubtree -> UnfoldKeysTreeNodeInfo -> LoadedKeysSubtree
collapseKeysTreeNode subtree nodeToCollapse =
  if List.isEmpty nodeToCollapse.path then
    {subtree | loadedNodes = List.map (\node -> 
      case node of
        UnfoldKeyTreeNode keyInfo -> 
          if keyInfo.name == nodeToCollapse.name then
            CollapsedKeyTreeNode <| CollapsedKeysTreeNodeInfo keyInfo.path keyInfo.name
          else
            node
        _ -> node
    ) subtree.loadedNodes}
  else
    {subtree | loadedNodes = List.map (\node ->
      case node of
        UnfoldKeyTreeNode keyInfo ->
          if Just keyInfo.name == List.head nodeToCollapse.path then
            UnfoldKeyTreeNode {keyInfo | loadedChildren = collapseKeysTreeNode keyInfo.loadedChildren {nodeToCollapse | path = List.drop 1 nodeToCollapse.path}}            
          else
            node
        _ -> node
    ) subtree.loadedNodes}



httpErrorToString : Http.Error -> String
httpErrorToString err = 
  case err of
    Http.BadUrl string -> "Wrong request url"
    Http.Timeout -> "Request timeout"
    Http.NetworkError -> "Network error"
    Http.BadStatus _ ->  "Got error response"
    Http.BadPayload _ _ ->  "Cannot decode response"


valueAddingErrorToString : Http.Error -> String
valueAddingErrorToString err =
 case err of
  Http.BadStatus response -> 
    if response.status.code == 409 then "Key already exists"
    else httpErrorToString err
  _ -> httpErrorToString err