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
import View.AboutModal as AboutModal
import View.Helpers exposing (..)
import Dialog

view : Model -> Html Msg
view model =
  div [ class "container" ] [
      navbar model,
      maybeDrawWorkspace model,
      Dialog.view (
        if model.addKeyModalShown then
          AddKeyModal.view model
        else if model.aboutWindowShown then
          AboutModal.view model
        else
          Nothing
      )
  ]

navbar : Model -> Html Msg
navbar model = 
  div [] [
    nav [ class "navbar navbar-default" ] [
        div [ class "container-fluid" ] [
            div [ class "navbar-header" ] [
              a [id "logo", class "navbar-brand", href "#", onClick LogoClick] [ 
                img [src "static/img/logo.png", alt "Radish"] []
              ]
            ],
            ul [ class "navbar-right nav navbar-nav" ] [
              li [ class "dropdown" ] [
                a [ class "dropdown-toggle", attribute "data-toggle" "dropdown" ] [
                    i [class "fa fa-server"] [],
                    text " ",
                    drawChosenServerName model,
                    span [class "caret"] []
                ],
                ul [ class "dropdown-menu" ] 
                  <| List.map (\server -> 
                      li [] [
                        a [href "#", onClick (ChosenServer server.name)] [
                          text <| server.name ++ " ",                          
                          drawIfFalse server.connectionCheckPassed 
                            <| span [class "fa fa-warning text-warning", title "Connection problems"] []
                        ]  
                      ] 
                    ) <| Dict.values model.loadedServers.servers
              ],
              li [ class <| 
                  "dropdown " ++ 
                    ((\chosenServer -> if model.chosenServer == Nothing then "hidden" else "") model.chosenServer) 
                ] [
                a [ class "dropdown-toggle", attribute "data-toggle" "dropdown" ] [
                  i [class "fa fa-database"] [],
                  text (" db # " ++ (toString model.chosenDatabaseNum)),
                  span [class "caret"] []
                ],
                ul [class "dropdown-menu"]
                  <| List.map (\dbNum -> 
                    li [] [
                      a [href "#", onClick (DatabaseChosen dbNum)] [
                        text <| toString dbNum
                      ]
                    ]
                  ) <| List.range 0 
                    <| Maybe.withDefault 0 
                    <| Maybe.andThen (\server -> Just (server.databasesCount - 1 )) (getChosenServer model)
              ]
            ]
      ]
    ]
  ]

drawChosenServerName : Model -> Html Msg
drawChosenServerName model =
  case model.chosenServer of
    Just serverName -> text serverName
    Nothing -> text "Choose server"


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
          drawKeysPanel model
      ],
      div [class "col-md-8"] [
          valuesPanel model 
      ]
    ]

drawKeysPanel : Model -> Html Msg
drawKeysPanel model =
  case model.chosenKeysViewType of 
    KeysListView ->
      drawListKeysPanel model
    KeysTreeView -> 
      drawTreeKeysPanel model

drawListKeysPanel : Model -> Html Msg
drawListKeysPanel model =
  div [class "panel panel-default"] [
    div [class "panel-heading"] [
        Html.form [class "form-inline"] [
          div [class "input-group form-group"] [
            input [type_ "text", class "form-control", placeholder "*", value model.keysMask, onInput KeysMaskChanged] [],
            label [class "input-group-addon"] [i [class "fa fa-search"] []]
          ],
          div [class "pull-right"] [
            drawKeysActionsButton model
          ]
        ]   
    ],
    div [class "panel-body keys-list", attribute "style" (keysListHeightAttribute model)] [
        div [class "pull-right label label-default"] [
          text "Found ",
          text <| toString model.loadedKeys.foundKeysCount,
          text " keys"
        ],
        div [class "clearfix"] [],
        div [class "top-buffer list-group"]
          (List.map (\key ->
            a [
              class ("list-group-item break-word " ++ (if model.chosenKey == Just key then "active" else "")), 
              onClick (KeyChosen key)
              ] [text key]
          ) model.loadedKeys.keys)
    ], 
    div [class "panel-footer"] [
        drawPager model.loadedKeys.pagesCount model.loadedKeys.currentPage KeysPageChanged
    ]
  ]
  
drawTreeKeysPanel : Model -> Html Msg
drawTreeKeysPanel model =
  div [class "panel panel-default"] [
    div [class "panel-heading"] [
        div [class "pull-right"] [
          drawKeysActionsButton model
        ],
        div [class "clearfix"] []
    ],
    div [class "panel-body keys-list", attribute "style" (keysListHeightAttribute model)] [
      div [] [
        div [class "pull-right"] [
          drawKeysActionsButton model
        ]
      ],
      div [class "list-group"] 
        (drawKeysTreeBranch model model.loadedKeysTree)
    ]
  ]

keysListHeightAttribute : Model -> String
keysListHeightAttribute model =
  "height:" ++ (toString <| toFloat model.windowSize.height * 0.75) ++ "px;"

drawKeysActionsButton : Model -> Html Msg
drawKeysActionsButton model =
  div [class "btn-group"] [
    button [class "btn btn-default dropdown-toggle", attribute "data-toggle" "dropdown"] [
      text "Action",
      span [class "caret"] []
    ],
    ul [class "dropdown-menu"] [
      li [] [
        a [onClick KeysListUpdateInitialized] [
          i [class "fa fa-refresh"] [],
          text " Refresh keys"
        ]
      ],
      li [] [
        a [onClick ShowAddKeyModal] [
            i [class "fa fa-plus-square"] [],
            text " Add key"
        ]
      ],
      li [] [
        drawToggleKeysViewButton model
      ]
    ]
  ]


drawToggleKeysViewButton : Model -> Html Msg
drawToggleKeysViewButton model = 
  case model.chosenKeysViewType of 
    KeysListView ->
      a [onClick TreeKeyViewChosen] [
          i [class "fa fa-folder-o"] [],
          text " Switch to tree keys view"
      ]
    KeysTreeView -> 
      a [onClick ListKeyViewChosen] [
          i [class "fa fa-list"] [],
          text " Switch to list keys view"
      ]

drawKeysTreeBranch : Model -> LoadedKeysSubtree -> List (Html Msg)
drawKeysTreeBranch model loadedSubtree =
  List.map (drawKeysTreeNode model) loadedSubtree.loadedNodes


drawKeysTreeNode : Model -> KeysTreeNode -> Html Msg
drawKeysTreeNode model node =
  case node of
    CollapsedKeyTreeNode keyInfo ->
      a [class "list-group-item break-word", onClick <| KeysTreeCollapsedNodeClick keyInfo] [
        span [class "fa fa-plus"] [],
        text (" " ++ keyInfo.name)
      ]
    KeysTreeLeaf keyInfo ->
      a [
        class ("list-group-item break-word " ++ (if model.chosenKey == Just keyInfo.key then "active" else "")), 
        onClick (KeyChosen keyInfo.key)
      ] [text keyInfo.name]
    UnfoldKeyTreeNode keyInfo ->
      span [class "list-group-item break-word"] ([
        span [class "fa fa-minus", onClick <| KeysTreeUnfoldNodeClick keyInfo] [text <| " " ++ keyInfo.name]
      ] ++ (drawKeysTreeBranch model keyInfo.loadedChildren))