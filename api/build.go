package api

import (
    "github.com/abemedia/push-deploy/agent"
    "github.com/abemedia/push-deploy/models"
    "net/http"
    "encoding/json"
    "github.com/abemedia/push-deploy/lib/session"
    "github.com/gorilla/mux"
)

func Builds(w http.ResponseWriter, r *http.Request) {
    if _, ok := getProjectAuthorised(w,r); !ok {
        return
    }
    builds, err := models.Builds.List(mux.Vars(r)["project_id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(builds)
}

func BuildNew(w http.ResponseWriter, r *http.Request) {
    p, err := models.Projects.Get(mux.Vars(r)["project_id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    s, err := session.Get(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    u := s.Values["user"].(models.User)
    if !u.Admin && u.ID != p.UserID {
        http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
        return
    }
    
    h := map[string]string{
        // todo: use head commit message
        "message": "Manual build (from Push Deploy)",
        "author": u.Name,
        "email": u.Email,
    }
    
    go agent.Run(h, p)
}

func BuildView(w http.ResponseWriter, r *http.Request) {
    if _, ok := getProjectAuthorised(w,r); !ok {
        return
    }
    b, err := models.Builds.Get(mux.Vars(r)["build_id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(b)
    
}

func BuildDelete(w http.ResponseWriter, r *http.Request) {
    if _, ok := getProjectAuthorised(w,r); !ok {
        return
    }
    
    err := models.Builds.Delete(mux.Vars(r)["build_id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
