module Command.Keys exposing (getKeysPage, deleteKey, addKey)

import Model.Model exposing (Model, RedisKey, LoadedKeys, keyTypeFromAlias)
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