module Command.Keys exposing (getKeysPage, deleteKey, addKey)

import Model.Model exposing (Model, RedisKey, LoadedKeys, KeysTreeNode(..), LoadedKeysSubtree, UnfoldKeysTreeNodeInfo, CollapsedKeysTreeNodeInfo, keyTypeFromAlias)
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


getKeysSubtree : Model -> String -> Int -> Cmd Msg
getKeysSubtree model prefix offset =
    case model.chosenServer of
        Just chosenServer ->
            let
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys-tree/" ++ (Http.encodeUri prefix) ++ "?offset=" ++ (toString offset)
            in
                Http.send KeysTreeSubtreeLoaded (Http.get url keysSubtreeDecoder)
        Nothing ->
            Cmd.none

keysSubtreeDecoder : Decode.Decoder LoadedKeysSubtree
keysSubtreeDecoder = 
    Decode.map3 LoadedKeysSubtree (Decode.field "KeysCount" Decode.int) (Decode.field "Path" <| Decode.list Decode.string) (Decode.field "Nodes" decodeKeysSubtreeNodes)

decodeKeysSubtreeNodes : Decode.Decoder (List KeysTreeNode)
decodeKeysSubtreeNodes = 
    Decode.list 
        <| (Decode.field "HasChildren" Decode.bool |> Decode.andThen decodeKeysTreeNode)

decodeKeysTreeNode : Bool -> Decode.Decoder KeysTreeNode
decodeKeysTreeNode hasChildren =
    if hasChildren then
        collapsedKeysTreeNodeDecoder
    else
        leafKeysTreeNodeDecoder


collapsedKeysTreeNodeDecoder : Decode.Decoder KeysTreeNode
collapsedKeysTreeNodeDecoder =
    Decode.map CollapsedKeyTreeNode <| Decode.map CollapsedKeysTreeNodeInfo (Decode.field "Name" Decode.string)

leafKeysTreeNodeDecoder : Decode.Decoder KeysTreeNode
leafKeysTreeNodeDecoder =
    Decode.map KeysTreeLeaf (Decode.field "Name" Decode.string)

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