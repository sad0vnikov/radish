module Update.Update exposing (Msg, update)

import Model.Model exposing (..)

type Msg = NoOp | ChosenServer String

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    ChosenServer server ->
      ({model | chosenServer = Just server}, Cmd.none)
    _ ->
      (model, Cmd.none)




