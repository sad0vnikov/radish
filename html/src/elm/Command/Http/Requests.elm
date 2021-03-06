module Command.Http.Requests exposing (put, delete, post)

import Http

post : String -> Http.Body -> Http.Request ()
post url body = 
    Http.request
    { method = "POST"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
    }

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


delete : String -> Http.Request ()
delete url = 
    Http.request
    { method = "DELETE"
        , headers = []
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
    }