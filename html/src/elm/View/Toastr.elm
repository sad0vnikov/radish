port module View.Toastr exposing (..)


port toastError : String -> Cmd msg
port toastInfo : String -> Cmd msg
port toastWarning : String -> Cmd msg
port toastSuccess : String -> Cmd msg