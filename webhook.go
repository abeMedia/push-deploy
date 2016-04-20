package main

import (
    "strings"
    "encoding/json"
    "net/http"
    "bytes"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/abemedia/push-deploy/agent"
    "github.com/abemedia/push-deploy/models"
)


func init() {
    var router = mux.NewRouter()
    router.HandleFunc("/hook/{id:[0-9]+}", webhook)
    http.Handle("/hook/", router)
}

func webhook(w http.ResponseWriter, r *http.Request) {
    project, err := models.Projects.Get(mux.Vars(r)["id"])
    if err != nil {
        fmt.Fprintln(w, "Error:", err.Error())
        return
    }
    
    var hook map[string]string
    
    // detect webhook type
    if strings.HasPrefix(r.Header.Get("User-Agent"), "GitHub") {
        buf := new(bytes.Buffer)
        buf.ReadFrom(r.Body)
        hook = webhookGithub(buf.Bytes())
    } else {
        // parse json into struct
        decoder := json.NewDecoder(r.Body)
        err = decoder.Decode(&hook)
    }
    
    if err != nil {
        fmt.Fprintln(w, "Error:", err.Error())
        return
    }
    
    // check push is coming from correct branch
    if hook["branch"] != project.Branch {
        fmt.Fprintln(w, "Error: Wrong branch")
        return
    }
    
    // queue the build in separate goroutine
    go agent.Run(hook, project)
    //fmt.Printf("%s", project)
    fmt.Fprint(w, "Hook received")
}

func webhookGithub(payload []byte) map[string]string {
    var data interface{}
    json.Unmarshal(payload, &data)
    hook := data.(map[string]interface{})
    ref := strings.Split(hook["ref"].(string), "/")
    return map[string]string{
        //"repo": "github.com/" + hook["repository"].(map[string]interface{})["full_name"].(string),
        "branch": ref[len(ref)-1],
        "message": hook["head_commit"].(map[string]interface{})["message"].(string),
        "author": hook["head_commit"].(map[string]interface{})["author"].(map[string]interface{})["name"].(string),
        "email": hook["head_commit"].(map[string]interface{})["author"].(map[string]interface{})["email"].(string),
    }
}