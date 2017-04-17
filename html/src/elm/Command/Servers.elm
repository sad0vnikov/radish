module Command.Servers exposing (getServersList)

import Update.Msg exposing (Msg(..))
import Json.Decode as Decode exposing (..)
import Model.Model exposing (Server,Model)
import Http
import Dict exposing (..)

getServersList : Model -> Cmd Msg
getServersList model = 
    let 
        url = model.api.url ++ "/servers"
    in
        Http.send ServersListLoaded (Http.get url serversListDecoder)


serversListDecoder : Decode.Decoder (Dict String Server)
serversListDecoder =
    Decode.dict <| Decode.map3 Server (Decode.field "Name" Decode.string) (Decode.field "Host" Decode.string) (Decode.field "Port" Decode.int)
