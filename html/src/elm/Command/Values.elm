module Command.Values exposing (getKeyValues, deleteValue, updateValue )

import Update.Msg exposing (Msg(KeyValuesLoaded, ValueDeleted, ValueUpdated))
import Json.Decode as Decode
import Json.Encode as Encode
import Model.Model exposing (Model, LoadedKeys, LoadedValues, KeyType, RedisKey,
  LoadedValues(..), RedisValuesPage, RedisValue, KeyType(..), RedisValues(..), 
  StringRedisValue, ListRedisValue, ZSetRedisValue, getChosenServerAndKey, getLoadedKeyType)
import Http
import Command.Http.Requests exposing (put, delete)
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

deleteValue : Model -> String -> Cmd Msg
deleteValue model valueToDelete =
    case getChosenServerAndKey model of
        Just (chosenServer,chosenKey) ->
            let
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys/" ++ valueDeleteUrlPath chosenKey (getLoadedKeyType model.loadedValues) valueToDelete
            in
                Http.send ValueDeleted (Http.getString url)
        Nothing ->
            Cmd.none
    
valueDeleteUrlPath :  RedisKey -> KeyType -> String -> String
valueDeleteUrlPath chosenKey keyType value =
    case keyType of
        StringRedisKey -> ""
        HashRedisKey -> "/hashes/" ++ chosenKey ++ "/values/" ++ value ++ "/delete"
        SetRedisKey -> "/sets/" ++ chosenKey ++ "/values/" ++ value ++ "/delete"
        ZSetRedisKey -> "/zsets/" ++ chosenKey ++ "/values/" ++ value ++ "/delele"
        ListRedisKey -> "/lists/" ++ chosenKey ++ "/values/" ++ value ++ "/delete"
        UnknownRedisKey -> ""

updateValue : Model -> String -> String -> Cmd Msg
updateValue model value newValue = 
    case getChosenServerAndKey model of
        Just (chosenServer,chosenKey) ->
            let
                url = model.api.url ++ "/servers/" ++ chosenServer ++ "/keys" ++ valueEditUrlPath chosenKey (getLoadedKeyType model.loadedValues) value
            in
                Http.send ValueUpdated (put url (Http.jsonBody <| valueEditJsonRequest model chosenKey))
        Nothing ->
            Cmd.none

valueEditJsonRequest : Model -> RedisKey -> Encode.Value
valueEditJsonRequest model redisKey =
    (Encode.object [("Value", Encode.string model.editingValueToSave)])


valueEditUrlPath : RedisKey -> KeyType -> String  -> String
valueEditUrlPath chosenKey keyType value  =
    case keyType of
        StringRedisKey -> "/strings/" ++ chosenKey
        HashRedisKey -> "/hashes/" ++ chosenKey ++ "/values/" ++ value ++ "/update"
        SetRedisKey -> "/sets/" ++ chosenKey ++ "/values/" ++ value ++ "/update"
        ZSetRedisKey -> "/zsets/" ++ chosenKey ++ "/values/" ++ value ++ "/update"
        ListRedisKey -> "/lists/" ++ chosenKey ++ "/values/" ++ value ++ "/update"
        UnknownRedisKey -> ""

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
        , decodeHashValues
        , decodeSetValues
    ]


decodeListValues : Decode.Decoder RedisValues
decodeListValues =
    Decode.map ListRedisValues <|
        Decode.list <| Decode.map2 ListRedisValue (Decode.field "Index" Decode.int) (Decode.field "Value" Decode.string)

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
            (Decode.field "Member" Decode.string) 