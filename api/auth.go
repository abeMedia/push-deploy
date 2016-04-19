package api

import (
    //"fmt"
    "net/http"
    "github.com/abemedia/push-deploy/models"
    "github.com/abemedia/push-deploy/lib/session"
)


func auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	    
	    // todo: remove gorilla sessions and use cookies from http package
	    // fmt.Println(r.Cookies())
	    // http.SetCookie(w, &http.Cookie{Name: "csrftoken",Value:"abcd",Expires:expiration})
        
        // check if user is already logged in
        s, err := session.Get(r)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        
        if user, ok := s.Values["user"].(models.User); ok && user.ID > 0 {
            h(w, r)
            return
        }
        if username, password, ok := r.BasicAuth(); ok {
            if user, err := models.Users.Auth(username, password); err == nil && user.ID > 0 {
                // create user session 
                s.Values["user"] = user
                //r.User = user
                s.Save(r, w)
                // Delegate request to the given handle
                h(w, r)
                return
            }
        }

        // Request Basic Authentication otherwise
        //w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
        http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
    }
}