module Main exposing (..)

import Html exposing (program)

import View.ServerPage exposing (..)
import Model.Model exposing (..)
import Update.Update exposing (update)
import Update.Msg exposing (Msg(..))
import Flags exposing (Flags)
import Command.Servers exposing (getServersList)
import View.ConfirmationDialog exposing (..)

main : Program Flags Model Msg
main =
  Html.programWithFlags { init = init, view = view, update = Update.Update.update, subscriptions = subscriptions }


subscriptions : Model -> Sub Msg
subscriptions model = 
  dialogClosed (\s -> 
      case s of
        "ok" ->
          UserConfirmation
        _ ->
          UserConfirmationCancel
    )


init : Flags -> (Model, Cmd Msg)
init flags =
  let
    model = initModel flags
  in
    (model, getServersList model)
