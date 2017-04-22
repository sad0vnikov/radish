module Command.Values exposing (getKeyValues)

import Update.Msg exposing (Msg(KeyValuesLoaded))
import Json.Decode as Decode
import Model.Model exposing (Model, LoadedKeys, LoadedValues, KeyType, 
  LoadedValues(..), RedisValuesPage, RedisValue, KeyType(..), RedisValues(..), StringRedisValue, ZSetRedisValue, getChosenServerAndKey)
import Http
import Maybe exposing (andThen)

getKeyValues : Model -> Cmd Msg
getKeyValues model =
    case getChosenServerAndKey model of
        Just (chosenServer,chosenKey) ->
            let
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys/" ++ chosenKey ++ "/values?page=" ++ (toString model.loadedKeys.currentPage)
            in
                Http.send KeyValuesLoaded (Http.get url valuesDecoder)
        Nothing ->
            Cmd.none


    

valuesDecoder : Decode.Decoder LoadedValues
valuesDecoder =
    Decode.oneOf [
        decodeRedisValuesPage
        , decodeRedisSingleValue
    ]

decodeRedisSingleValue : Decode.Decoder LoadedValues
decodeRedisSingleValue =
    Decode.map SingleRedisValue <| Decode.map2 RedisValue
        (Decode.field "Value" Decode.string)
        (Decode.field "KeyType" Decode.string |> Decode.andThen decodeKeyType)
        

decodeRedisValuesPage : Decode.Decoder LoadedValues
decodeRedisValuesPage =
    Decode.map MultipleRedisValues <| Decode.map4 RedisValuesPage
        (Decode.field "Values" decodeValuesList)     
        (Decode.field "PagesCount" Decode.int)
        (Decode.field "PageNum" Decode.int)        
        (Decode.field "KeyType" Decode.string |> Decode.andThen decodeKeyType) 

decodeKeyType : String -> Decode.Decoder KeyType
decodeKeyType status =
    Decode.succeed
        (
            case status of
                "string" -> StringRedisKey
                "hash" -> HashRedisKey
                "set" -> SetRedisKey
                "zset" -> ZSetRedisKey
                "list" -> ListRedisKey
                _ -> UnknownRedisKey
        )

decodeValuesList : Decode.Decoder RedisValues
decodeValuesList =
    Decode.oneOf [
        decodeListValues
        , decodeZSetValues
    ]


decodeListValues : Decode.Decoder RedisValues
decodeListValues =
    Decode.map ListRedisValues <|
        Decode.list Decode.string

decodeHashValues : Decode.Decoder RedisValues
decodeHashValues = 
    Decode.map HashRedisValues <| Decode.dict Decode.string

decodeSetValues : Decode.Decoder RedisValues
decodeSetValues =    
    Decode.map SetRedisValues <| Decode.list Decode.string

decodeZSetValues : Decode.Decoder RedisValues
decodeZSetValues =
    Decode.map ZSetRedisValues <| Decode.list <|
        Decode.map2 ZSetRedisValue
            (Decode.field "Score" Decode.int)
            (Decode.field "Value" Decode.string) 