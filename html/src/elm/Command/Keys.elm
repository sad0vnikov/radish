module Command.Keys exposing (getKeysPage, deleteKey)

import Model.Model exposing (Model, RedisKey, LoadedKeys)
import Update.Msg exposing (Msg(..))
import Http
import Json.Decode as Decode

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


deleteKey : String -> Model -> Cmd Msg
deleteKey key model =
    case model.chosenServer of
        Just chosenServer ->
            let
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys/" ++ key ++ "/delete"
            in
                Http.send KeyDeleted (Http.getString url)
        Nothing ->
            Cmd.none