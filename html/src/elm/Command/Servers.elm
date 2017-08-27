module Command.Servers exposing (getServersList)

import Update.Msg exposing (Msg(..))
import Json.Decode as Decode exposing (..)
import Model.Model exposing (Server,Model, ServerKeyspaceStat, ServerStat)
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
    Decode.dict <| Decode.map7 Server 
        (Decode.field "Name" Decode.string) 
        (Decode.field "Host" Decode.string) 
        (Decode.field "Port" Decode.int)
        (Decode.field "DatabasesCount" Decode.int)
        (Decode.field "ConnectionCheckPassed" Decode.bool)
        (Decode.field "ServerStat" <| decodeServerStat)
        (Decode.field "KeyspaceStat" decodeServerKeyspaceStat)

decodeServerStat : Decode.Decoder (Maybe ServerStat)
decodeServerStat =
    Decode.nullable <| Decode.map7 ServerStat
        (Decode.field "ConnectedClientsCount" Decode.int)
        (Decode.field "RedisVersion" Decode.string)
        (Decode.field "UptimeInSeconds" Decode.int)
        (Decode.field "UsedMemoryHuman" Decode.string)
        (Decode.field "UsedMemoryBytes" Decode.int)
        (Decode.field "MaxMemoryHuman" Decode.string)
        (Decode.field "MaxMemoryBytes" Decode.int)

decodeServerKeyspaceStat : Decode.Decoder (Maybe (Dict String ServerKeyspaceStat))
decodeServerKeyspaceStat =
    Decode.nullable <| Decode.dict <| Decode.map ServerKeyspaceStat (Decode.field "KeysCount" Decode.int)