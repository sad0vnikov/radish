module Command.Keys exposing (getKeysPage)

import Model.Model exposing (Model, RedisKey)
import Update.Msg exposing (Msg(..))
import Http
import Json.Decode as Decode

getKeysPage : Model -> Cmd Msg
getKeysPage model =
    case model.chosenServer of
        Just chosenServer ->
            let
                keysMask = sanitizeKeysMask model.keysMask
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys/bymask/" ++ keysMask
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

keysListDecoder : Decode.Decoder (List RedisKey)
keysListDecoder =
    Decode.at ["Keys"] (Decode.list Decode.string)
