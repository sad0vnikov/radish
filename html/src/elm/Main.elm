module Main exposing (..)

import Html exposing (program)

import View.ServerPage exposing (..)
import Model.Model exposing (..)
import Update.Update exposing (update, Msg)
import Flags exposing (Flags)
import Command.Servers exposing (getServersList)

main : Program Flags Model Msg
main =
  Html.programWithFlags { init = init, view = view, update = Update.Update.update, subscriptions = subscriptions }


subscriptions : Model -> Sub Msg
subscriptions model = Sub.none


init : Flags -> (Model, Cmd Msg)
init flags =
  let
    model = initModel flags
  in
    (model, getServersList model)
