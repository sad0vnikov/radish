module View.Pagination exposing (drawPager)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import List exposing (range)
import Update.Msg exposing (Msg)


drawPager : Int -> Int -> (Int -> Msg) -> Html Msg
drawPager pagesCount curPage onClickMsg = 
    nav [] [
        ul [class "pagination"] (
            (range 1 pagesCount) 
                |> List.map (\pageNum -> 
                    li (markPageIfActive pageNum curPage) [
                        a [href "#", onClick (onClickMsg pageNum)] [text (toString pageNum)]
                    ]
                )
        )
    ]

markPageIfActive : Int -> Int -> List (Html.Attribute msg)
markPageIfActive pageNum chosenPage =
    if pageNum /= chosenPage then
        []
    else
        [class "active"]