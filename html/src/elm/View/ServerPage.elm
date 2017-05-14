module View.ServerPage exposing (view)

import Html exposing (..)
import Html.Events exposing (onClick, onInput)

import Model.Model exposing (..)
import Html.Attributes exposing (..)
import Dict exposing (..)
import List exposing (range)
import Update.Msg exposing (Msg(..))
import View.Pagination exposing (drawPager)
import View.Values exposing (..)
import View.AddKeyModal as AddKeyModal

view : Model -> Html Msg
view model =
  div [ class "container" ] [   
      navbar model,
      maybeDrawWorkspace model,
      AddKeyModal.view model
  ]

navbar : Model -> Html Msg
navbar model = 
  div [] [
    nav [ class "navbar navbar-default" ] [
        div [ class "container-fluid" ] [
            div [ class "navbar-header" ] [
              a [class "navbar-brand" ] [ text "Radish" ]
            ],
            ul [ class "navbar-right nav navbar-nav" ] [
              li [ class "dropdown" ] [
                a [ class "dropdown-toggle", attribute "data-toggle" "dropdown" ] [
                    i [class "fa fa-server"] [],
                    text " Choose server"
                ],
                ul [ class "dropdown-menu" ] 
                  <| List.map (\server -> 
                      li [] [
                        a [href "#", onClick (ChosenServer server.name)] [text server.name]
                      ] 
                    ) <| Dict.values model.loadedServers.servers
              ]
          
          ]
      ]
    ]
  ]

maybeDrawWorkspace : Model -> Html Msg
maybeDrawWorkspace model =
  case model.chosenServer of
    Just server ->
      workspace model
    Nothing ->
      div [] []

workspace : Model -> Html Msg
workspace model = 
  div [class "row" ] [
      div [class "col-md-4"] [
        div [class "panel panel-default"] [
          div [class "panel-heading"] [
              
          ],
          div [class "panel-body"] [
              Html.form [class "form-inline"] [
                div [class "input-group form-group"] [
                  input [type_ "text", class "form-control", placeholder "*", value model.keysMask, onInput KeysMaskChanged] [],
                  label [class "input-group-addon"] [i [class "fa fa-search"] []]
                ], 
                button [class "btn btn-success pull-right", onClick ShowAddKeyModal] [
                    i [class "fa fa-plus-square"] [],
                    text " Add key"
                ]
              ],
              div [class "list-group"]
                (List.map (\key ->
                  a [
                    class ("list-group-item " ++ (if model.chosenKey == Just key then "active" else "")), 
                    onClick (KeyChosen key)
                    ] [text key]
                ) model.loadedKeys.keys)        
          ], 
          div [class "panel-footer"] [
              drawPager model.loadedKeys.pagesCount model.loadedKeys.currentPage KeysPageChanged
          ]
        ]
      ],
      div [class "col-md-8"] [
          valuesPanel model 
      ]
    ]