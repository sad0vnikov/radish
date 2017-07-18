module View.Values exposing (..)

import Dict exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Update.Msg exposing (Msg(..))
import Html.Events exposing (onClick, onInput)
import View.Pagination exposing (drawPager)
import View.Helpers exposing (drawIfFalse)


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
      h2 [class "break-word"] [
        text <| chosenKey ++ " ",
        small [] [drawKeyType model.loadedValues],
        button [class "btn btn-sm btn-danger pull-right", onClick (KeyDeletionConfirm chosenKey)] [
          i [class "fa fa-remove"] [],
          text " Delete key"
        ]
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
        ListRedisValues values -> listKeyValues model values redisValuesPage.pagesCount redisValuesPage.currentPage
        HashRedisValues values -> hashKeyValues model values redisValuesPage.pagesCount redisValuesPage.currentPage
        SetRedisValues values -> setKeyValues model values redisValuesPage.pagesCount redisValuesPage.currentPage
        ZSetRedisValues values -> sortedSetValues model values redisValuesPage.pagesCount redisValuesPage.currentPage
        _ ->
            div [] []


drawSingleRedisValue : Model -> String -> RedisValue -> Html Msg
drawSingleRedisValue model key redisValue =
    stringKeyValues model key redisValue

drawIconIfValueIsBinary : Bool -> Html Msg
drawIconIfValueIsBinary isBinary =
  if isBinary then
    span [] [
      span [class "fa fa-file-image-o"] [],
      text " "
    ]
  else
    span [] []
 
stringKeyValues : Model -> String -> RedisValue -> Html Msg
stringKeyValues model key value =
  div [] [
    div [class "btn-toolbar"] [
      drawIfFalse value.isBinary <|
        button [class "btn btn-sm", onClick <| ValueToEditSelected (value.value, value.value)] [
            i [class "fa fa-edit"] [],
            text " Edit"
        ]
    ],
    div [class "top-buffer"] [
        drawValueOrEditField model value
    ]
  ]



drawValueOrEditField : Model -> RedisValue -> Html Msg
drawValueOrEditField model value =
  case model.editingValue of
    Nothing ->
      div [class "well break-word"] [
          drawIconIfValueIsBinary value.isBinary,
          text value.value
      ]
    Just valueToEdit ->
      div[class "row"] [
        div [class "col-md-10"] [
           textarea [class "form-control", onInput EditedValueChanged] [text model.editingValueToSave]
        ],
        div [class "col-md-2"] [
          button [class "btn btn-sm btn-primary btn-block", onClick <| ValueUpdateInitialized value.value] [
            i [class "fa fa-save"] []
          ],
          button [class "btn btn-sm btn-default btn-block", onClick ValueEditingCanceled] [
             i [class "fa fa-times"] []
          ]           
        ]
      ]

hashKeyValues : Model -> (Dict String HashRedisValue) -> Int -> Int -> Html Msg
hashKeyValues model values pagesCount currentPage = 
 div [] [
   table [class "table hash-values"] [
     thead [] [
       th [class "key"] [text "key"],
       th [class "value"] [text "value"],
       th [class "buttons"] []
     ],
     tbody [] <| (Dict.values <| Dict.map (drawHashValueOrEditFields model) values) ++ (List.singleton <| maybeDrawHashValueAddFields model)
   ],
   div [class "row"] [
     div [class "col-xs-8"] [
        drawPager pagesCount currentPage ValuesPageChanged
     ],
     div [class "col-xs-4"] [
        button [class "btn btn-sm btn-primary add-new-value-btn", onClick AddingValueStart] [
          i [class "fa fa-plus"] [],
          text " Add a new value"
        ]
     ]
   ]
 ]

drawHashValueOrEditFields : Model -> String -> HashRedisValue -> Html Msg
drawHashValueOrEditFields model hashKey hashValue =
  case model.editingValue of
    Nothing -> drawHashValueRow model hashKey hashValue.value hashValue.isBinary
    Just (key, valueToEdit) -> 
      if valueToEdit == hashKey then
        drawHashValueEditFields model hashKey hashValue.value
      else
        drawHashValueRow model hashKey hashValue.value hashValue.isBinary

drawHashValueRow : Model -> String -> StringRedisValue -> Bool -> Html Msg
drawHashValueRow model hashKey hashValue isBinary = 
    tr [] [
         td [class "break-word key"] [
            drawIconIfValueIsBinary isBinary,
            text hashKey
         ], 
         td [class "break-word value"] [
            text hashValue
         ],
         td [class "buttons"] [
           drawIfFalse isBinary <|
            button [class "btn btn-sm btn-edit", onClick <| HashValueToEditSelected (hashKey, hashValue)] [
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
    td [class "key"] [
      input [class "form-control", Html.Attributes.value model.editingHashKeyToSave, onInput EditedHashKeyChanged] []
    ],
    td [class "value"] [
      input [class "form-control", Html.Attributes.value model.editingValueToSave, onInput EditedValueChanged] []
    ],
    td [class "buttons"] [
      button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized hashKey] [
        i [class "fa fa-save"] []
      ],
      button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
        i [class "fa fa-times"] []
      ]
    ]
  ]

maybeDrawHashValueAddFields : Model -> Html Msg
maybeDrawHashValueAddFields model =
  if model.isAddingValue then
    drawHashValueAddFields model
  else
    tr [] []

drawHashValueAddFields : Model -> Html Msg
drawHashValueAddFields model = 
  tr [] [
    td [class "key"] [
      input [class "form-control", Html.Attributes.value model.addingHashKey, onInput AddingHashKeyChanged] []
    ],
    td [class "value"] [
      input [class "form-control", Html.Attributes.value model.addingValue, onInput AddingValueChanged] []
    ],
    td [class "buttons"] [
      button [class "btn btn-sm btn-primary", onClick AddingValueInitialized] [
        i [class "fa fa-save"] []
      ],
      button [class "btn btn-sm btn-default", onClick AddingValueCancel] [
        i [class "fa fa-times"] []
      ]
    ]
  ]

listKeyValues : Model -> (List ListRedisValue) -> Int -> Int -> Html Msg
listKeyValues model values pagesCount currentPage = 
  div [] [
    table [class "table list-values"] [
      thead [] [
        th [class "index"] [text "index"],
        th [class "value"] [text "value"],
        th [class "buttons"] []
      ],
      tbody [] <| (List.map (drawListValueOrEditFields model) values) ++ (List.singleton <| maybeDrawListValueAddFields model)
    ],
    div [class "row"] [
      div [class "col-xs-8"] [
        drawPager pagesCount currentPage ValuesPageChanged
      ],
      div [class "col-xs-4"] [
        button [class "btn btn-sm btn-primary add-new-value-btn", onClick AddingValueStart] [
          i [class "fa fa-plus"] [],
          text " Add a new list member"
        ]
      ]
      
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
        td [class "break-word index"] [
            text <| toString listMember.index
        ],
        td [class "break-word value"] [
            drawIconIfValueIsBinary listMember.isBinary,
            text listMember.value
        ],
        td [class "buttons"] [
            drawIfFalse listMember.isBinary <|
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
      td [class "index"] [
        text <| toString listMember.index
      ],
      td [class "value"] [
        input [class "form-control", onInput EditedValueChanged, Html.Attributes.value model.editingValueToSave] []
      ],
      td [class "buttons"] [
        button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized (toString listMember.index)] [
          i [class "fa fa-save"] []
        ],
        button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
          i [class "fa fa-times"] []
        ]
      ]
  ]

maybeDrawListValueAddFields : Model -> Html Msg
maybeDrawListValueAddFields model =
  if model.isAddingValue then
    drawListValueAddFields model
  else 
    tr [] []


drawListValueAddFields : Model -> Html Msg
drawListValueAddFields model = 
    tr [] [
        td [class "index"] [
          
        ],
        td [class "value"] [
          input [class "form-control", onInput AddingValueChanged, Html.Attributes.value model.addingValue] []
        ],
        td [class "buttons"] [
          button [class "btn btn-sm btn-primary", onClick <| AddingValueInitialized] [
            i [class "fa fa-save"] []
          ],
          button [class "btn btn-sm btn-default", onClick AddingValueCancel] [
            i [class "fa fa-times"] []
          ]
        ]
    ]
  

setKeyValues : Model -> (List SetRedisValue) -> Int -> Int -> Html Msg
setKeyValues model values pagesCount currentPage = 
  div [] [
    table [class "table set-values"] [
      thead [] [

        th [class "value"] [text "value"],
        th [class "buttons"] []
      ],
      tbody [] <| (List.map (drawSetValueRowOrEditFields model) values) ++ (List.singleton <| maybeDrawSetValueAddFields model)
    ],
    div [class "row"] [
      div [class "col-xs-8"] [
        drawPager pagesCount currentPage ValuesPageChanged
      ],
      div [class "col-xs-4"] [
        button [class "btn btn-sm btn-primary add-new-value-button", onClick AddingValueStart] [
          i [class "fa fa-plus"] [],
          text " Add a new set member"
        ]
      ]
      
    ]
  ]

drawSetValueRowOrEditFields : Model -> SetRedisValue -> Html Msg
drawSetValueRowOrEditFields model value =
  case model.editingValue of
    Just (key, valueReference) ->
      if valueReference == value.value then
        drawSetEditFieldsRow model value.value
      else
        drawSetValueRow model value.value value.isBinary

    Nothing ->
      drawSetValueRow model value.value value.isBinary


drawSetValueRow : Model -> StringRedisValue -> Bool-> Html Msg
drawSetValueRow model value isBinary =
    tr [] [
        td [class "break-word value"] [
            drawIconIfValueIsBinary isBinary,
            text value
          ],
        td [class "buttons"] [
            drawIfFalse isBinary <|
              button [class "btn btn-sm", onClick <| ValueToEditSelected (value, value)] [i [class "fa fa-pencil"] []],
            button [class "btn btn-sm btn-danger", onClick (ValueDeletionConfirm value)] [i [class "fa fa-remove"] []]
        ]
    ]

drawSetEditFieldsRow : Model -> StringRedisValue -> Html Msg
drawSetEditFieldsRow model value =
    tr [] [
      td [class "value"] [
        input [class "form-control", onInput EditedValueChanged, Html.Attributes.value model.editingValueToSave] []
      ],
      td [class "buttons"] [
        button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized value] [
          i [class "fa fa-save"] []
        ],
        button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
          i [class "fa fa-times"] []
        ]        
      ]
    ]

maybeDrawSetValueAddFields : Model -> Html Msg
maybeDrawSetValueAddFields model =
  if model.isAddingValue then
    drawSetValueAddFields model
  else
    tr [] []

drawSetValueAddFields : Model -> Html Msg
drawSetValueAddFields model =
  tr [] [
      td [class "value"] [
        input [class "form-control", onInput AddingValueChanged, Html.Attributes.value model.addingValue] []
      ],
      td [class "buttons"] [
        button [class "btn btn-sm btn-primary", onClick AddingValueInitialized] [
          i [class "fa fa-save"] []
        ],
        button [class "btn btn-sm btn-default", onClick AddingValueCancel] [
          i [class "fa fa-times"] []
        ]        
      ]
    ]

sortedSetValues: Model -> (List ZSetRedisValue) -> Int -> Int -> Html Msg
sortedSetValues model values pagesCount currentPage = 
  div [] [
    table [class "table zset-values"] [
      thead [] [
        th [class "score"] [text "score"],
        th [class "value"] [text "value"],
        th [class "buttons"] []
      ],
      tbody [] <| (List.map (drawSortedSetValueRowOrEditFields model) values) ++ (List.singleton <| maybeDrawSortedSetAddFields model)
    ],
    div [class "row"] [
      div [class "col-xs-8"] [
        drawPager pagesCount currentPage ValuesPageChanged
      ], 
      div [class "col-xs-4"] [
        button [class "btn btn-sm btn-primary add-new-value-button", onClick AddingValueStart] [
          i [class "fa fa-plus"] [],
          text " Add a new member"
        ]
      ]
      
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
        td [class "break-word score"] [
            text <| toString value.score
        ],
        td [class "break-word value"] [
            drawIconIfValueIsBinary value.isBinary,
            text value.value
        ],
        td [class "buttons"] [
            drawIfFalse value.isBinary <|
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
      td [class "score"] [
          input [class "form-control", Html.Attributes.value <| toString model.editingScoreToSave, onInput EditedScoreChanged] []
      ],
      td [class "value"] [
          input [class "form-control", Html.Attributes.value model.editingValueToSave, onInput EditedValueChanged] []          
      ],
      td [class "buttons"] [
          button [class "btn btn-sm btn-primary", onClick <| ValueUpdateInitialized value.value] [
            i [class "fa fa-save"] []
          ],
          button [class "btn btn-sm btn-default", onClick ValueEditingCanceled] [
            i [class "fa fa-times"] []
          ]
      ]
    ]

maybeDrawSortedSetAddFields : Model -> Html Msg
maybeDrawSortedSetAddFields model =
  if model.isAddingValue then
    drawSortedSetAddFields model
  else
    tr [] []

drawSortedSetAddFields : Model -> Html Msg
drawSortedSetAddFields model =
  tr [] [
      td [class "score"] [
          input [class "form-control", Html.Attributes.value <| toString model.addingZSetScore, onInput AddingZSetScoreChanged] []
      ],
      td [class "values"] [
          input [class "form-control", Html.Attributes.value model.addingValue, onInput AddingValueChanged] []          
      ],
      td [class "buttons"] [
          button [class "btn btn-sm btn-primary", onClick <| AddingValueInitialized] [
            i [class "fa fa-save"] []
          ],
          button [class "btn btn-sm btn-default", onClick AddingValueCancel] [
            i [class "fa fa-times"] []
          ]
      ]
  ]