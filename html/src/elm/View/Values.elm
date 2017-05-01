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
        MultipleRedisValues values -> drawMultipleRedisValues model key values
        SingleRedisValue value ->  drawSingleRedisValue model key value


drawMultipleRedisValues : Model -> String -> RedisValuesPage -> Html Msg
drawMultipleRedisValues model key redisValuesPage =
    case redisValuesPage.values of
        ListRedisValues values -> listKeyValues model values
        HashRedisValues values -> hashKeyValues model values
        SetRedisValues values -> setKeyValues model values
        ZSetRedisValues values -> sortedSetValues model values
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
      button [class "btn btn-sm", onClick <| ValueToEditSelected (value, value)] [
          i [class "fa fa-edit"] [],
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
           button [class "btn btn-sm btn-primary btn-block", onClick <| ValueUpdateInitialized value] [
             i [class "fa fa-save"] []
           ],
           button [class "btn btn-sm btn-default btn-block", onClick ValueEditingCanceled] [
             i [class "fa fa-times"] []
           ]           
        ]
      ]

hashKeyValues : Model -> (Dict String StringRedisValue) -> Html Msg
hashKeyValues model values = 
 div [] [
   table [class "table"] [
     thead [] [
       th [] [text "key"],
       th [] [text "value"],
       th [] []
     ],
     tbody [] <| Dict.values <| Dict.map (drawHashValueOrEditFields model) values
   ]
 ]

drawHashValueOrEditFields : Model -> String -> StringRedisValue -> Html Msg
drawHashValueOrEditFields model hashKey hashValue =
  case model.editingValue of
    Nothing -> drawHashValueRow model hashKey hashValue
    Just (key, valueToEdit) -> 
      if valueToEdit == hashKey then
        drawHashValueEditFields model hashKey hashValue
      else
        drawHashValueRow model hashKey hashValue

drawHashValueRow : Model -> String -> StringRedisValue -> Html Msg
drawHashValueRow model hashKey hashValue = 
    tr [] [
         td [] [
            text hashKey
         ], 
         td [] [
            text hashValue
         ],
         td [] [
           button [class "btn btn-sm btn-edit", onClick <| ValueToEditSelected (hashKey, hashValue)] [
             i [class "fa fa-pencil"] []
           ],
           button [class "btn btn-sm btn-danger", onClick (ValueDeletionConfirm hashKey)] [
             i [class "fa fa-remove"] []
           ]
         ]
    ]

drawHashValueEditFields : Model -> String -> StringRedisValue -> Html Msg
drawHashValueEditFields model hashKey hashValue =
  tr [] [
    td [] [
      text hashKey
    ],
    td [] [
      input [class "form-control", Html.Attributes.value model.editingValueToSave, onInput EditedValueChanged] []
    ],
    td [] [
      button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized hashKey] [
        i [class "fa fa-save"] []
      ],
      button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
        i [class "fa fa-times"] []
      ]
    ]
  ]


listKeyValues : Model -> (List ListRedisValue) -> Html Msg
listKeyValues model values = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "index"],
        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map (drawListValueOrEditFields model) values
    ]
  ]

drawListValueOrEditFields : Model -> ListRedisValue -> Html Msg
drawListValueOrEditFields model value =
  case model.editingValue of
    Just (key, valueReference) ->
      if valueReference == toString value.index then
        drawListValueEditFields model value
      else 
        drawListValueRow model value 

    Nothing -> drawListValueRow model value    

drawListValueRow : Model -> ListRedisValue -> Html Msg
drawListValueRow  model listMember = 
    tr [] [
        td [] [
            text <| toString listMember.index
        ],
        td [] [
            text listMember.value
        ],
        td [] [
            button [class "btn btn-sm", onClick <| ValueToEditSelected (toString listMember.index, listMember.value)] [
              i [class "fa fa-pencil"] []
            ],
            button [class "btn btn-sm btn-danger", onClick (ValueDeletionConfirm <| toString listMember.index)] [
              i [class "fa fa-remove"] []
            ]
        ]
    ]

drawListValueEditFields : Model -> ListRedisValue -> Html Msg
drawListValueEditFields model listMember =
  tr [] [
    td [] [
        text <| toString listMember.index
      ],
      td [] [
        input [class "form-control", onInput EditedValueChanged, Html.Attributes.value model.editingValueToSave] []
      ],
      td [] [
        button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized (toString listMember.index)] [
          i [class "fa fa-save"] []
        ],
        button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
          i [class "fa fa-times"] []
        ]
      ]
  ]

setKeyValues : Model -> (List StringRedisValue) -> Html Msg
setKeyValues model values = 
  div [] [
    table [class "table"] [
      thead [] [

        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map (drawSetValueRowOrEditFields model) values
    ]
  ]

drawSetValueRowOrEditFields : Model -> StringRedisValue -> Html Msg
drawSetValueRowOrEditFields model value =
  case model.editingValue of
    Just (key, valueReference) ->
      if valueReference == value then
        drawSetEditFieldsRow model value
      else
        drawSetValueRow model value

    Nothing ->
      drawSetValueRow model value


drawSetValueRow : Model -> StringRedisValue -> Html Msg
drawSetValueRow model value =
    tr [] [
        td [] [
            text value
          ],
        td [] [
            button [class "btn btn-sm", onClick <| ValueToEditSelected (value, value)] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger", onClick (ValueDeletionConfirm value)] [i [class "fa fa-remove"] []]
        ]
    ]

drawSetEditFieldsRow : Model -> StringRedisValue -> Html Msg
drawSetEditFieldsRow model value =
    tr [] [
      td [] [
        input [class "form-control", onInput EditedValueChanged, Html.Attributes.value model.editingValueToSave] []
      ],
      td [] [
        button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized value] [
          i [class "fa fa-save"] []
        ],
        button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
          i [class "fa fa-times"] []
        ]        
      ]
    ]

sortedSetValues: Model -> (List ZSetRedisValue) -> Html Msg
sortedSetValues model values = 
  div [] [
    table [class "table"] [
      thead [] [
        th [] [text "score"],
        th [] [text "value"],
        th [] []
      ],
      tbody [] <| List.map (drawSortedSetValueRowOrEditFields model) values
    ]
  ]

drawSortedSetValueRowOrEditFields : Model -> ZSetRedisValue -> Html Msg
drawSortedSetValueRowOrEditFields model value = 
  case model.editingValue of
    Just (key, valueReference) ->
      if valueReference == value.value then
        drawSortedSetEditFields model value
      else
        drawSortedSetValueRow model value
    Nothing ->
      drawSortedSetValueRow model value

drawSortedSetValueRow : Model -> ZSetRedisValue -> Html Msg
drawSortedSetValueRow model value =
     tr [] [
        td [] [
            text <| toString value.score
        ],
        td [] [
            text value.value
        ],
        td [] [
            button [class "btn btn-sm btn-edit", onClick <| ZSetValueToEditSelected (value.value, value.value, value.score)] [
              i [class "fa fa-pencil"] []
            ],
            button [class "btn btn-sm btn-danger", onClick (ValueDeletionConfirm value.value)] [
              i [class "fa fa-remove"] []
            ]
        ]
    ]

drawSortedSetEditFields : Model -> ZSetRedisValue -> Html Msg
drawSortedSetEditFields model value =
    tr [] [
      td [] [
          input [class "form-control", Html.Attributes.value <| toString model.editingScoreToSave, onInput EditedScoreChanged] []
      ],
      td [] [
          input [class "form-control", Html.Attributes.value model.editingValueToSave, onInput EditedValueChanged] []          
      ],
      td [] [
          button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized value.value] [
            i [class "fa fa-save"] []
          ],
          button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
            i [class "fa fa-times"] []
          ]
      ]
    ]