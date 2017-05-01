module View.Values exposing (..)

import Dict exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Update.Msg exposing (Msg(..))
import Html.Events exposing (onClick, onInput)


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
        SingleRedisValue value ->  drawSingleRedisValue model key value


drawMultipleRedisValues : String -> RedisValuesPage -> Html Msg
drawMultipleRedisValues key redisValuesPage =
    case redisValuesPage.values of
        ListRedisValues values -> listKeyValues values
        HashRedisValues values -> hashKeyValues values
        SetRedisValues values -> setKeyValues values
        ZSetRedisValues values -> sortedSetValues values
        _ ->
            div [] []


drawSingleRedisValue : Model -> String -> RedisValue -> Html Msg
drawSingleRedisValue model key redisValue =
    stringKeyValues model key redisValue.value

 
stringKeyValues : Model -> String -> StringRedisValue -> Html Msg
stringKeyValues model key value =
  div [] [
    div [class "btn-toolbar"] [
      button [class "btn btn-sm btn-danger", onClick (KeyDeletionConfirm key)] [
          i [class "fa fa-remove"] [],
          text " Delete"
      ],
      button [class "btn btn-sm"] [
          i [class "fa fa-edit", onClick <| ValueToEditSelected value] [],
          text " Edit"
      ]
    ],
    div [class "top-buffer"] [
        drawValueOrEditField model value
    ]
  ]

drawValueOrEditField : Model -> StringRedisValue -> Html Msg
drawValueOrEditField model value =
  case model.editingValue of
    Nothing ->
      div [class "well"] [
          text value
      ]
    Just valueToEdit ->
      div[class "row"] [
        div [class "col-md-10"] [
           textarea [class "form-control", onInput EditedValueChanged] [text model.editingValueToSave]
        ],
        div [class "col-md-2"] [
           button [class "btn btn-sm btn-primary btn-block"] [i [class "fa fa-save", onClick <| ValueUpdateInitialized value] []],
           button [class "btn btn-sm btn-default btn-block", onClick ValueEditingCanceled] [i [class "fa fa-times"] []]           
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
           button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove", onClick (ValueDeletionConfirm key)] []]
         ]
    ]

listKeyValues : (List ListRedisValue) -> Html Msg
listKeyValues values = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "index"],
        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map drawListValueRow values
    ]
  ]

drawListValueRow : ListRedisValue -> Html Msg
drawListValueRow listMember = 
    tr [] [
        td [] [
            text <| toString listMember.index
        ],
        td [] [
            text listMember.value
        ],
        td [] [
            button [class "btn btn-sm btn-edit"] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove", onClick (ValueDeletionConfirm <| toString listMember.index)] []]
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
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove", onClick (ValueDeletionConfirm value)] []]
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
            button [class "btn btn-sm btn-danger"] [i [class "fa fa-remove", onClick (ValueDeletionConfirm value.value)] []]
        ]
    ]