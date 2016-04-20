package session

import (
    "net/http"
    "github.com/gorilla/sessions"
)

var store *sessions.CookieStore // var store = sessions.NewCookieStore([]byte())

func init() {
    cs := &sessions.CookieStore{
        Options: &sessions.Options{
            Path:   "/",
            MaxAge: 86400 * 30,
            Secure: false,
        },
    }

    cs.MaxAge(cs.Options.MaxAge)
    store = cs  
}

func Get(r *http.Request) (*sessions.Session, error) {
    return store.Get(r, "session")
}
