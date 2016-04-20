package api

import (
    "path"
    "fmt"
    "io/ioutil"
    "net/http"
    . "github.com/abemedia/push-deploy/lib/config"
    "github.com/gorilla/mux"
)

func Log(w http.ResponseWriter, r *http.Request) {
    // open log file
    file, err := ioutil.ReadFile(path.Join(Config.CachePath, mux.Vars(r)["id"], "logs/build.log"))
    if err != nil {
        fmt.Fprintln(w, "No Log found")
    }
    fmt.Fprintln(w, string(file))
}
