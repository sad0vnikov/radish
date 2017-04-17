module View.ServerPage exposing (view)

import Html exposing (..)
import Html.Events exposing (onClick, onInput)
import Model.Model exposing (..)
import Html.Attributes exposing (..)
import Dict exposing (..)
import Update.Msg exposing (Msg(..))

view : Model -> Html Msg
view model =
  div [ class "container" ] [   
      navbar model,
      maybeDrawWorkspace model
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
              Html.form [] [
                div [class "input-group"] [
                  input [type_ "text", class "form-control", placeholder "*", value model.keysMask, onInput KeysMaskChanged] [],
                  label [class "input-group-addon"] [i [class "fa fa-search"] []]
                ]
              ],
              ul [class "list-group"]
                (List.map (\key ->
                  li [class "list-group-item"] [text key]
                ) model.loadedKeys)        
          ]
        ]
      ],
      div [class "col-md-8"] [
          actionPanel model 
      ]
    ]
    

actionPanel : Model -> Html Msg
actionPanel model =
  div [class "panel panel-default"] [
    div [class "panel-heading"] [

    ],
    div [class "panel-body"] [
      div [] [
        ul [class "nav nav-tabs"] [
          li [class "active"] [
            a [href "#"] [
              i [class "fa fa-edit"] [],
              text " Key editor"
            ]
          ]
        ]
      ],
      div [] [
        maybeDrawKeyEditor model
      ]
    ]
  ]


maybeDrawKeyEditor : Model -> Html Msg
maybeDrawKeyEditor model =
  case model.chosenKey of
    Just key ->
      keyEditor model
    Nothing ->
      div [] []

keyEditor : Model -> Html Msg
keyEditor model = 
  div [] [
    div [class "page-header"] [
      h2 [] [
        text "a:1 ",
        small [] [text "hash value"]
      ]   
    ],
    listKeyValue model
  ]
 
stringKeyValue : Model -> Html Msg
stringKeyValue model =
  div [] [
    div [class "btn-toolbar"] [
      button [class "btn btn-sm btn-danger"] [
          i [class "fa fa-remove"] [],
          text " Delete"
      ]
    ],
    div [class "top-buffer"] [
        div [class "well"] [
          text "this is the key value"
        ]
    ]
  ]

hashKeyValue : Model -> Html Msg
hashKeyValue model = 
 div [] [
   table [class "table"] [
     thead [] [
       th [] [text "key"],
       th [] [text "value"],
       th [] []
     ],
     tbody [] [
       tr [] [
         td [] [
            text "bb"
         ], 
         td [] [
            text "2"
         ],
         td [] [
           button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
           button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
         ]
       ]
     ]
   ]
 ]

listKeyValue : Model -> Html Msg
listKeyValue model = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "index"],
        th [] [text "value"],
        th [] []
      ],
      tbody [] [
        tr [] [
          td [] [
            text "0"
          ], 
          td [] [
            text "a"
          ],
          td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
          ]
        ]
      ]
    ]
  ]

setKeyValue : Model -> Html Msg
setKeyValue model = 
  div [] [
    table [class "table"] [
      thead [] [

        th [] [text "value"],
        th [] []
      ],
      tbody [] [
        tr [] [
          td [] [
            text "some value"
          ],
          td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
          ]
        ]
      ]
    ]
  ]


sortedSetKeyValue : Model -> Html Msg
sortedSetKeyValue model = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "score"],
        th [] [text "value"],
        th [] []
      ],
      tbody [] [
        tr [] [
          td [] [
            text "1122"
          ],
          td [] [
            text "some value"
          ],
          td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
          ]
        ]
      ]
    ]
  ]