module Main exposing (..)

import Html.Events exposing ( onClick )
import Dict exposing (..)

import Html exposing (program)

import View.ServerPage exposing (..)
import Model.Model exposing (..)
import Update.Update exposing (update, Msg)


type alias Flags = {
  apiUrl: String
}

main : Program Flags Model Msg
main =
  Html.programWithFlags { init = init, view = view, update = Update.Update.update, subscriptions = subscriptions }


subscriptions : Model -> Sub Msg
subscriptions model = Sub.none


init : Flags -> (Model, Cmd Msg)
init flags =
  let
    servers =  Dict.empty 
      |> (Dict.insert "server1" (Server "server1" "127.0.0.1" 6379))
      |> (Dict.insert "server2" (Server "server2" "127.0.0.1" 6380))
  in
    ({
      api = { url = flags.apiUrl},
      loadedData = {
        servers = Dict.empty,
        loadedKeys = []
      }, 
      chosenServer = Maybe.Nothing,
      chosenKey = Maybe.Nothing
    }, Cmd.none)
