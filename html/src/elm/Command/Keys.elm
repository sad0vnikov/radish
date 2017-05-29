module Command.Keys exposing (getKeysPage, getKeysSubtree, deleteKey, addKey)

import Model.Model exposing (Model, RedisKey, LoadedKeys, KeysTreeNode(..), LoadedKeysSubtree, UnfoldKeysTreeNodeInfo, CollapsedKeysTreeNodeInfo, KeysTreeLeafInfo, keyTypeFromAlias)
import Update.Msg exposing (Msg(..))
import Http
import Json.Decode as Decode
import Command.Http.Requests exposing (delete)
import Command.Values exposing (addValue)

getKeysPage : Model -> Cmd Msg
getKeysPage model =
    case model.chosenServer of
        Just chosenServer ->
            let
                keysMask = sanitizeKeysMask model.keysMask
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys?mask=" ++ (Http.encodeUri keysMask) ++ "&page=" ++ toString model.loadedKeys.currentPage
            in
                Http.send KeysPageLoaded (Http.get url keysListDecoder)
        Nothing ->
            Cmd.none

sanitizeKeysMask : String -> String
sanitizeKeysMask mask = 
    if (String.trim mask) == "" then
        "*"
    else
        mask

keysListDecoder : Decode.Decoder LoadedKeys
keysListDecoder =
    Decode.map3 LoadedKeys (Decode.field "Keys" (Decode.list Decode.string)) (Decode.field "PagesCount" Decode.int) (Decode.field "Page" Decode.int)


getKeysSubtree : Model -> List String -> Cmd Msg
getKeysSubtree model path =
    case model.chosenServer of
        Just chosenServer ->
            let
                treePath = String.join "/" path
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys-tree?path=" ++ (Http.encodeUri treePath)
            in
                Http.send KeysTreeSubtreeLoaded (Http.get url keysSubtreeDecoder)
        Nothing ->
            Cmd.none

keysSubtreeDecoder : Decode.Decoder LoadedKeysSubtree
keysSubtreeDecoder = 
    Decode.map2 LoadedKeysSubtree (decodeSubtreePath) (decodeSubtreePath |> Decode.andThen decodeKeysSubtreeNodes)

decodeSubtreePath : Decode.Decoder (List String)
decodeSubtreePath = 
    Decode.field "Path" <| Decode.list Decode.string

decodeKeysSubtreeNodes : List String -> Decode.Decoder (List KeysTreeNode)
decodeKeysSubtreeNodes path = 
    Decode.field "Nodes" <| Decode.list
        <| (Decode.field "HasChildren" Decode.bool |> Decode.andThen (decodeKeysTreeNode path))

decodeKeysTreeNode : List String -> Bool -> Decode.Decoder KeysTreeNode
decodeKeysTreeNode path hasChildren =
    if hasChildren then
        collapsedKeysTreeNodeDecoder path
    else
        leafKeysTreeNodeDecoder


collapsedKeysTreeNodeDecoder : List String -> Decode.Decoder KeysTreeNode
collapsedKeysTreeNodeDecoder path =
    Decode.map CollapsedKeyTreeNode <| Decode.map2 CollapsedKeysTreeNodeInfo (Decode.succeed path) (Decode.field "Name" Decode.string)

leafKeysTreeNodeDecoder : Decode.Decoder KeysTreeNode
leafKeysTreeNodeDecoder =
    Decode.map KeysTreeLeaf <| Decode.map2 KeysTreeLeafInfo  (Decode.field "Key" Decode.string) (Decode.field "Name" Decode.string)

addKey : Model -> Cmd Msg
addKey model =
    case model.chosenServer of
        Just serverName ->
            addValue serverName model.keyToAddName (keyTypeFromAlias model.keyToAddType) {model | chosenKey = Just model.keyToAddName, addingValue = "0", addingHashKey = "0", addingZSetScore = 0}
        Nothing ->
            Cmd.none


deleteKey : String -> Model -> Cmd Msg
deleteKey key model =
    case model.chosenServer of
        Just chosenServer ->
            let
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys/" ++ key
            in
                Http.send KeyDeleted (delete url)
        Nothing ->
            Cmd.none