module View.AboutModal exposing (view)

import Model.Model exposing (Model, availableKeyTypes, keyTypeName, keyTypeAlias, KeyType)
import Update.Msg exposing (Msg(..))
import Html exposing (..)
import Html.Attributes exposing (..)
import Dialog

view : Model -> Maybe (Dialog.Config Msg)
view model =
    Just {
        closeMessage = Just AboutWindowClose,
        containerClass = Nothing,
        header = Just (text <| "Radish " ++ model.appInfo.version),
        body = Just (
            div [] [
                p [] [
                    strong [] [
                        text "Author: "
                    ],
                    text "Alexander Sadovnikov"
                ],
                p [] [
                    a [href "http://github.com/sad0vnikov/radish"] [
                        i [class "fa fa-2x fa-github"] []
                    ]
                ]
            ]
        ),
        footer = Nothing
    }


