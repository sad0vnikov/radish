module View.Values exposing (..)

import Dict exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Update.Msg exposing (Msg(..))
import Html.Events exposing (onClick)


import Model.Model exposing (..)

valuesPanel : Model -> Html Msg
valuesPanel model =
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
        maybeDrawValuesEditor model
      ]
    ]
  ]

maybeDrawValuesEditor : Model -> Html Msg
maybeDrawValuesEditor model =
  case model.chosenKey of
    Just key ->
      keyEditor model key
    Nothing ->
      div [] []

keyEditor : Model -> String -> Html Msg
keyEditor model chosenKey = 
  div [] [
    div [class "page-header"] [
      h2 [] [
        text <| chosenKey ++ " ",
        small [] [drawKeyType model.loadedValues]
      ]   
    ],
    drawKeyValuesEditorByType chosenKey model
  ]

drawKeyType : LoadedValues -> Html Msg
drawKeyType loadedValues =
  case getLoadedKeyType loadedValues of
    StringRedisKey -> text "string"
    ListRedisKey -> text "list"
    SetRedisKey -> text "set"
    ZSetRedisKey -> text "zset"
    HashRedisKey -> text "hash"
    UnknownRedisKey -> text "unknown key type"

drawKeyValuesEditorByType : String -> Model -> Html Msg
drawKeyValuesEditorByType key model = 
    case model.loadedValues of
        MultipleRedisValues values -> drawMultipleRedisValues key values
        SingleRedisValue value ->  drawSingleRedisValue key value


drawMultipleRedisValues : String -> RedisValuesPage -> Html Msg
drawMultipleRedisValues key redisValuesPage =
    case redisValuesPage.values of
        ListRedisValues values -> listKeyValues values
        HashRedisValues values -> hashKeyValues values
        SetRedisValues values -> setKeyValues values
        ZSetRedisValues values -> sortedSetValues values
        _ ->
            div [] []


drawSingleRedisValue : String -> RedisValue -> Html Msg
drawSingleRedisValue key redisValue =
    stringKeyValues key redisValue.value

 
stringKeyValues : String -> StringRedisValue -> Html Msg
stringKeyValues key value =
  div [] [
    div [class "btn-toolbar"] [
      button [class "btn btn-sm btn-danger", onClick (KeyDeletionConfirm key)] [
          i [class "fa fa-remove"] [],
          text " Delete"
      ]
    ],
    div [class "top-buffer"] [
        div [class "well"] [
          text value
        ]
    ]
  ]

hashKeyValues : (Dict String StringRedisValue) -> Html Msg
hashKeyValues values = 
 div [] [
   table [class "table"] [
     thead [] [
       th [] [text "key"],
       th [] [text "value"],
       th [] []
     ],
     tbody [] <| Dict.values <| Dict.map drawHashValueRow values
   ]
 ]

drawHashValueRow : String -> StringRedisValue -> Html Msg
drawHashValueRow key value = 
    tr [] [
         td [] [
            text key
         ], 
         td [] [
            text value
         ],
         td [] [
           button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
           button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
         ]
    ]

listKeyValues : (List StringRedisValue) -> Html Msg
listKeyValues values = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map drawListValueRow values
    ]
  ]

drawListValueRow : StringRedisValue -> Html Msg
drawListValueRow value = 
    tr [] [
        td [] [
            text value
        ],
        td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
        ]
    ]

setKeyValues : (List StringRedisValue) -> Html Msg
setKeyValues values = 
  div [] [
    table [class "table"] [
      thead [] [

        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map drawSetValueRow values
    ]
  ]

drawSetValueRow : StringRedisValue -> Html Msg
drawSetValueRow value =
    tr [] [
        td [] [
            text value
          ],
        td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
        ]
    ]



sortedSetValues: (List ZSetRedisValue) -> Html Msg
sortedSetValues values = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "score"],
        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map drawSortedSetValueRow values
    ]
  ]

drawSortedSetValueRow : ZSetRedisValue -> Html Msg
drawSortedSetValueRow value =
     tr [] [
        td [] [
            text <| toString value.score
        ],
        td [] [
            text value.value
        ],
        td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove"] []]
        ]
    ]