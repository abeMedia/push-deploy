package api

import (
    "net/http"
    "github.com/abemedia/push-deploy/lib/session"
    "github.com/abemedia/push-deploy/models"
    "encoding/json"
)

func Login(w http.ResponseWriter, r *http.Request) {
    s, err := session.Get(r)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    var body map[string]string
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&body)
    
    if body["username"] != "" && body["password"] != "" {
        if user, err := models.Users.Auth(body["username"], body["password"]); err == nil {
            s.Values["user"] = user
            s.Save(r, w)
            json.NewEncoder(w).Encode(user)
            return
        }
    }
    http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
