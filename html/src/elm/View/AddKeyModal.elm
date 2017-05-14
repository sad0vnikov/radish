module View.AddKeyModal exposing (view)

import Model.Model exposing (Model, availableKeyTypes, keyTypeName, keyTypeAlias, KeyType)
import Update.Msg exposing (Msg(..))
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput, on, targetValue)
import Dialog exposing (view)
import View.HtmlEvents.Events exposing (onChange)

view : Model -> Html Msg
view model =
    Dialog.view 
        (if model.addKeyModalShown then
            Just {
                closeMessage = Just CloseAddKeyModal,
                containerClass = Nothing,
                header = Just (text "Add a new key"),
                body = Just (
                    Html.form [] [
                        div [class "form-group"] [
                            label [for "add-key-type-selector"] [text "Key name"],                            
                            select [id "add-key-type-selector", class "form-control", onChange KeyToAddTypeChanged] 
                                <| List.map (keyTypeOption model) availableKeyTypes
                        ],
                        div [class "form-group"] [
                            label [for "add-key-name-selector"] [text "Key name"],
                            input [type_ "text", id "add-key-name-selector", class "form-control", onInput KeyToAddNameChanged] []
                        ]
                    ]
                ),
                footer = Just (
                    button [class "btn btn-primary", onClick AddNewKey] [text "Add key"]
                )
            }
        else 
            Nothing
        )


keyTypeOption : Model -> KeyType -> Html Msg
keyTypeOption model keyType =
    let 
        keyTypeStringName = keyTypeName keyType
    in
        option [value <| keyTypeAlias keyType, selected (model.keyToAddType == keyTypeStringName)] [
            text <| keyTypeStringName
        ]
