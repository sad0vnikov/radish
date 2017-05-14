module View.HtmlEvents.Events exposing (onChange)

import Html exposing (..)
import Html.Events exposing (..)
import Json.Decode as Json


onChange : (String -> msg) -> Attribute msg
onChange tagger =
    on "change" (Json.map tagger targetValue)