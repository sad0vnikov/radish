module Command.Values exposing (getKeyValues, deleteValue, updateValue, addValueForChosenKey, addValue)

import Update.Msg exposing (Msg(KeyValuesLoaded, ValueDeleted, ValueUpdated, ValueAdded))
import Json.Decode as Decode
import Json.Encode as Encode
import Model.Model exposing (Model, LoadedKeys, LoadedValues, KeyType, RedisKey,
  LoadedValues(..), RedisValuesPage, RedisValue, KeyType(..), RedisValues(..), 
  StringRedisValue, ListRedisValue, SetRedisValue, ZSetRedisValue, HashRedisValue, getChosenServerAndKey, getLoadedKeyType)
import Http
import Command.Http.Requests exposing (put, delete, post)
import Maybe exposing (andThen)

getKeyValues : Model -> Int -> Cmd Msg
getKeyValues model pageNum =
    case getChosenServerAndKey model of
        Just (chosenServer,chosenKey) ->
            let
                url = model.api.url 
                    ++ "/servers/" ++ chosenServer 
                    ++ "/keys/" ++ chosenKey 
                    ++ "/values?page=" ++ (toString pageNum)
                    ++ "&db=" ++ (toString model.chosenDatabaseNum)
                    ++ "&mask=" ++ (model.valuesMask)
            in
                Http.send KeyValuesLoaded (Http.get url valuesDecoder)
        Nothing ->
            Cmd.none

deleteValue : Model -> String -> Cmd Msg
deleteValue model valueToDelete =
    case getChosenServerAndKey model of
        Just (chosenServer,chosenKey) ->
            let
                url = model.api.url 
                    ++ "/servers/" ++ chosenServer 
                    ++ "/keys/" ++ valueDeleteUrlPath chosenKey (getLoadedKeyType model.loadedValues) valueToDelete
                    ++ "?db=" ++ toString model.chosenDatabaseNum
            in
                Http.send ValueDeleted (delete url)
        Nothing ->
            Cmd.none
    
valueDeleteUrlPath :  RedisKey -> KeyType -> String -> String
valueDeleteUrlPath chosenKey keyType value =
    case keyType of
        StringRedisKey -> ""
        HashRedisKey -> "/hashes/" ++ chosenKey ++ "/values/" ++ value
        SetRedisKey -> "/sets/" ++ chosenKey ++ "/values/" ++ value
        ZSetRedisKey -> "/zsets/" ++ chosenKey ++ "/values/" ++ value
        ListRedisKey -> "/lists/" ++ chosenKey ++ "/values/" ++ value
        UnknownRedisKey -> ""

addValueForChosenKey : Model -> Cmd Msg
addValueForChosenKey model =
    case getChosenServerAndKey model of 
        Just (chosenServer,chosenKey) ->
            let
                keyType = getLoadedKeyType model.loadedValues
            in
                addValue chosenServer chosenKey keyType model
        Nothing ->
            Cmd.none

addValue : String -> RedisKey -> KeyType -> Model -> Cmd Msg
addValue serverName keyName keyType model =
    let 
        url = model.api.url 
            ++ "/servers/" ++ serverName 
            ++ "/keys" ++ valueAddUrlPath keyName keyType
            ++ "?db=" ++ toString model.chosenDatabaseNum
    in
        Http.send ValueAdded (post url (Http.jsonBody <| valueAddJsonRequest model keyType))

valueAddUrlPath : RedisKey -> KeyType -> String
valueAddUrlPath chosenKey keyType =
    case keyType of
        ListRedisKey -> "/lists/" ++ chosenKey ++ "/values"
        HashRedisKey -> "/hashes/" ++ chosenKey ++  "/values"
        SetRedisKey -> "/sets/" ++ chosenKey ++ "/values"
        ZSetRedisKey -> "/zsets/" ++ chosenKey ++ "/values"        
        StringRedisKey -> "/strings/" ++ chosenKey
        _ -> ""

valueAddJsonRequest : Model -> KeyType -> Encode.Value
valueAddJsonRequest model keyType =
    case keyType of
        HashRedisKey -> (Encode.object [
                ("Value", Encode.string model.addingValue), 
                ("Key", Encode.string model.addingHashKey)
            ])
        ZSetRedisKey -> (Encode.object [
                ("Value", Encode.string model.addingValue),
                ("Score", Encode.int model.addingZSetScore)
            ])
        _ -> (Encode.object [("Value", Encode.string model.addingValue)])


updateValue : Model -> String -> String -> Cmd Msg
updateValue model value newValue = 
    case getChosenServerAndKey model of
        Just (chosenServer,chosenKey) ->
            let
                keyType = getLoadedKeyType model.loadedValues
                url = model.api.url 
                    ++ "/servers/" ++ chosenServer 
                    ++ "/keys" ++ valueEditUrlPath chosenKey keyType value
                    ++ "?db=" ++ toString model.chosenDatabaseNum
            in
                Http.send ValueUpdated (put url (Http.jsonBody <| valueEditJsonRequest model keyType))
        Nothing ->
            Cmd.none

valueEditJsonRequest : Model -> KeyType -> Encode.Value
valueEditJsonRequest model keyType =
    case keyType of
        ZSetRedisKey -> 
            (Encode.object [
                ("Value", Encode.string model.editingValueToSave),
                ("Score", Encode.int model.editingScoreToSave)
            ]) 
        HashRedisKey ->
            (Encode.object [
                ("Value", Encode.string model.editingValueToSave),
                ("NewHashKey", Encode.string model.editingHashKeyToSave)
            ])
        _ -> (Encode.object [("Value", Encode.string model.editingValueToSave)])   


valueEditUrlPath : RedisKey -> KeyType -> String  -> String
valueEditUrlPath chosenKey keyType value  =
    case keyType of
        StringRedisKey -> "/strings/" ++ chosenKey
        HashRedisKey -> "/hashes/" ++ chosenKey ++ "/values/" ++ value
        SetRedisKey -> "/sets/" ++ chosenKey ++ "/values/" ++ value
        ZSetRedisKey -> "/zsets/" ++ chosenKey ++ "/values/" ++ value
        ListRedisKey -> "/lists/" ++ chosenKey ++ "/values/" ++ value
        UnknownRedisKey -> ""

valuesDecoder : Decode.Decoder LoadedValues
valuesDecoder =
    Decode.oneOf [
        decodeRedisValuesPage
        , decodeRedisSingleValue
    ]

decodeRedisSingleValue : Decode.Decoder LoadedValues
decodeRedisSingleValue =
    Decode.map SingleRedisValue <| Decode.map3 RedisValue
        (Decode.field "Value" 
            <| Decode.field "Value" Decode.string)
        (Decode.field "KeyType" Decode.string |> Decode.andThen decodeKeyType)
        (Decode.field "Value"
            <| Decode.field "IsBinary" Decode.bool)
        

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
        Decode.list 
            <| Decode.map3 ListRedisValue 
                (Decode.field "Index" Decode.int) 
                (Decode.field "Value" <| Decode.field "Value" Decode.string) 
                (Decode.field "Value" <| Decode.field "IsBinary" Decode.bool)

decodeHashValues : Decode.Decoder RedisValues
decodeHashValues = 
    Decode.map HashRedisValues 
        <| Decode.dict
            (Decode.map2 HashRedisValue (Decode.field "Value" Decode.string) (Decode.field "IsBinary" Decode.bool))

decodeSetValues : Decode.Decoder RedisValues
decodeSetValues =    
    Decode.map SetRedisValues 
        <| Decode.list 
        <| Decode.map2 SetRedisValue (Decode.field "Value" Decode.string) (Decode.field "IsBinary" Decode.bool)

decodeZSetValues : Decode.Decoder RedisValues
decodeZSetValues =
    Decode.map ZSetRedisValues <| Decode.list <|
        Decode.map3 ZSetRedisValue
            (Decode.field "Score" Decode.int)
            (Decode.field "Member" <| Decode.field "Value" Decode.string)
            (Decode.field "Member" <| Decode.field "IsBinary" Decode.bool)
