package api

import (
    //"fmt"
    "github.com/abemedia/push-deploy/models"
    "net/http"
    "github.com/abemedia/push-deploy/lib/session"
    "encoding/json"
    "github.com/gorilla/mux"
)

func getUserAuthorised(w http.ResponseWriter, r *http.Request) (models.User, bool) {
    user, err := models.Users.Get(mux.Vars(r)["id"])
    s, err := session.Get(r)
    if err == nil {
        u := s.Values["user"].(models.User)
        if u.Admin || u.ID == user.ID {
            return user, true
        }
    }
    http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
    return user, false
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
    if _, ok := getUserAuthorised(w,r); !ok {
        return
    }
    var data models.User
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    id, err := models.Users.Add(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    u, err := models.Users.Get(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(u)
}

func UserView(w http.ResponseWriter, r *http.Request) {
    if u, ok := getUserAuthorised(w,r); ok {
        json.NewEncoder(w).Encode(u)
    }
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
    if _, ok := getUserAuthorised(w,r); !ok {
        return
    }
    var data models.User
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    id := mux.Vars(r)["id"]
    err = models.Users.Update(id, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    u, err := models.Users.Get(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(u)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
    if _, ok := getUserAuthorised(w,r); !ok {
        return
    }
    err := models.Users.Delete(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}


func Users(w http.ResponseWriter, r *http.Request) {
    if _, ok := getUserAuthorised(w,r); !ok {
        return
    }
    u, err := models.Users.List()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(u)
}