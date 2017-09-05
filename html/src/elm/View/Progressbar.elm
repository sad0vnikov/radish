module View.Progressbar exposing (drawProgressbar)

import Html exposing (..)
import Html.Attributes exposing (..)
import Update.Msg exposing (Msg)


drawProgressbar : Int -> Int -> String -> String -> Html Msg
drawProgressbar curValue maxValue curValueText maxValueText =
    let
        progressbarWidth = getProgressbarPercentage curValue maxValue
        progressBarClass = getProgressBarClass progressbarWidth
    in 
        div [class "progress"] [
            div [class <| "progress-bar progress-bar-" ++ progressBarClass, style [("width", toString progressbarWidth ++ "%")]] [
                span [class "text-muted"] [
                    text <| curValueText ++ "/" ++ maxValueText
                ]
            ]
        ]


getProgressbarPercentage : Int -> Int -> Float
getProgressbarPercentage curValue maxValue =
    toFloat curValue / toFloat maxValue * 100


getProgressBarClass : Float -> String
getProgressBarClass percentage =
    if percentage < 70 then
        "success"
    else if percentage < 90 then 
        "warning"
    else
        "danger"
    