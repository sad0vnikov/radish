module View.Helpers exposing (drawIfTrue, drawIfFalse)

import Html exposing (..)
import Update.Msg exposing (Msg(..))


drawIfTrue : Bool -> (Html Msg) -> Html Msg
drawIfTrue bool html =
    if bool then
        html
    else 
        span [] []

drawIfFalse : Bool -> (Html Msg) -> Html Msg
drawIfFalse bool html =
    if not bool then
        html
    else 
        span [] []