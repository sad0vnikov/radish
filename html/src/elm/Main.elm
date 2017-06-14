module Main exposing (..)

import Html exposing (program)

import View.ServerPage exposing (..)
import Model.Model exposing (..)
import Update.Update exposing (update)
import Update.Msg exposing (Msg(..))
import Flags exposing (Flags)
import Command.Servers exposing (getServersList)
import View.ConfirmationDialog exposing (..)
import Task
import Window


main : Program Flags Model Msg
main =
  Html.programWithFlags { init = init, view = view, update = Update.Update.update, subscriptions = subscriptions }


subscriptions : Model -> Sub Msg
subscriptions model = 
  Sub.batch [
    dialogClosed (\s -> 
      case s of
        "ok" ->
          UserConfirmation
        _ ->
          UserConfirmationCancel
    )
  ]


init : Flags -> (Model, Cmd Msg)
init flags =
  let
    model = initModel flags
  in
    (model, Cmd.batch [getServersList model, Task.perform windowSizeToMsg Window.size])

windowSizeToMsg : Window.Size -> Msg
windowSizeToMsg size = 
  WindowResized (size.width, size.height)