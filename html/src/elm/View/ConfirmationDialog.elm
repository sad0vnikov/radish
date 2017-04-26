port module View.ConfirmationDialog exposing (..)

port showConfirmationDialog : String -> Cmd msg
port dialogClosed : (String -> msg) -> Sub msg
