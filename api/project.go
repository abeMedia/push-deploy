package api

import (
    "github.com/abemedia/push-deploy/agent"
    "github.com/abemedia/push-deploy/models"
    "net/http"
    "encoding/json"
    "github.com/abemedia/push-deploy/lib/session"
    "github.com/gorilla/mux"
)

func getProjectAuthorised(w http.ResponseWriter, r *http.Request) (models.Project, bool) {
    project, err := models.Projects.Get(mux.Vars(r)["project_id"])
    s, err := session.Get(r)
    if err == nil {
        u := s.Values["user"].(models.User)
        if u.Admin || u.ID == project.UserID {
            return project, true
        }
    }
    http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
    return project, false
}

func ProjectAdd(w http.ResponseWriter, r *http.Request) {
    var d models.Project
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&d)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // set user_id to current user
    s, err := session.Get(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    d.UserID = s.Values["user"].(models.User).ID
    
    id, err := models.Projects.Add(d)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    p, err := models.Projects.Get(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(p)
}

func ProjectView(w http.ResponseWriter, r *http.Request) {
    if project, ok := getProjectAuthorised(w,r); ok {
        json.NewEncoder(w).Encode(project)
    }
}

func ProjectUpdate(w http.ResponseWriter, r *http.Request) {
    if _, ok := getProjectAuthorised(w,r); !ok {
        return
    }
    
    var data models.Project
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    id := mux.Vars(r)["project_id"]
    err = models.Projects.Update(id, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    p, err := models.Projects.Get(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(p)
}

func ProjectDelete(w http.ResponseWriter, r *http.Request) {
    if _, ok := getProjectAuthorised(w,r); !ok {
        return
    }
    
    err := models.Projects.Delete(mux.Vars(r)["project_id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func ProjectBuild(w http.ResponseWriter, r *http.Request) {
    d, err := models.Projects.Get(mux.Vars(r)["project_id"])
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
    if !u.Admin && u.ID != d.UserID {
        http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
        return
    }
    
    h := map[string]string{
        "message": "Manual build (from jekyll-deploy dashboard)",
        "author": u.Name,
        "email": u.Email,
    }
    
    go agent.Run(h, d)
}


func Projects(w http.ResponseWriter, r *http.Request) {
    s, err := session.Get(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    project, err := models.Projects.UserList(s.Values["user"].(models.User).ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(project)
}
