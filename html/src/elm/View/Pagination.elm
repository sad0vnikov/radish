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
                |> filterPageButtons pagesCount curPage
                |> addFirstAndLastPages pagesCount 
                |> List.map (\pageNum -> 
                    if pageNum > 0 then
                        li (markPageIfActive pageNum curPage) [
                            a [href "#", onClick (onClickMsg pageNum)] [text (toString pageNum)]
                        ]
                    else
                        li [class "disabled"] [
                            a [href "#"] [text "..."]
                        ]
                )
        )
    ]


addFirstAndLastPages : Int -> (List Int) -> (List Int)
addFirstAndLastPages lastPagesNum pages =
    (getFirstPageNumbers pages) ++ pages ++ (getLastPageNumbers lastPagesNum pages)

{-- a button for page 0 will be displayed as a button with ellipsis --}
getFirstPageNumbers : (List Int) -> (List Int)
getFirstPageNumbers pages =
    if List.head pages == Nothing then 
        []
    else if List.head pages == Just 2 then
        [1]
    else
        [1, 0]

{-- a button for page 0 will be displayed as a button with ellipsis --}
getLastPageNumbers : Int -> (List Int) -> (List Int)
getLastPageNumbers lastPageNum pages =
    let 
        reversedPages = List.reverse pages
    in
        if List.head reversedPages == Just lastPageNum then
            []
        else if List.head reversedPages == Just (lastPageNum - 1) then
            [lastPageNum]
        else if List.head reversedPages == Nothing then
            []
        else
            [0, lastPageNum]

filterPageButtons :  Int -> Int -> (List Int) -> (List Int)
filterPageButtons totalPagesCount currentPage pages =
    let
      pagerLinksToDrawCount = 7
      chosenPageNbhd = getChosenPageNeighbourhood totalPagesCount pagerLinksToDrawCount currentPage
    in
      List.filter (filterButton chosenPageNbhd) pages
    
filterButton : (Int, Int) -> Int -> Bool
filterButton  (chosenPageNbhdLeftEdge, chosenPageNbhdRightEdge) pageToFilterNum =
    if pageToFilterNum >= chosenPageNbhdLeftEdge && pageToFilterNum <= chosenPageNbhdRightEdge then
        True
    else
        False

getChosenPageNeighbourhood : Int -> Int -> Int -> (Int, Int)
getChosenPageNeighbourhood totalPagesCount pagesToDrawCount chosenPageNum  =
    let 
        curPageNbhdLeftHalfSize = (pagesToDrawCount - 2) // 2
        curPageNbhdRightHalfSize = pagesToDrawCount - curPageNbhdLeftHalfSize - 2
    in
        if totalPagesCount <= pagesToDrawCount then
            (2, pagesToDrawCount - 1)
        else if chosenPageNum - curPageNbhdLeftHalfSize <= 1 then
            (2, curPageNbhdRightHalfSize + curPageNbhdLeftHalfSize + 1)
        else if chosenPageNum + curPageNbhdRightHalfSize >= totalPagesCount then
            (totalPagesCount - curPageNbhdLeftHalfSize - curPageNbhdRightHalfSize, totalPagesCount - 1)
        else
            (chosenPageNum  - curPageNbhdLeftHalfSize, chosenPageNum + curPageNbhdRightHalfSize - 1)
        


markPageIfActive : Int -> Int -> List (Html.Attribute msg)
markPageIfActive pageNum chosenPage =
    if pageNum /= chosenPage then
        []
    else
        [class "active"]