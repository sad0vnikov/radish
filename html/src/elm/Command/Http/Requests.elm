module Command.Http.Requests exposing (put, delete)

import Http

put : String -> Http.Body -> Http.Request ()
put url body = 
    Http.request
    { method = "PUT"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
    }


delete : String -> Http.Body -> Http.Request ()
delete url body = 
    Http.request
    { method = "DELETE"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
    }