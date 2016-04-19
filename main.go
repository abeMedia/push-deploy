package main

import (
    "fmt"
    "net/http"
    "log"
    _ "github.com/abemedia/push-deploy/api"
    . "github.com/abemedia/push-deploy/lib/config"
    "github.com/gorilla/securecookie"
)

//var router = mux.NewRouter()
var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64),securecookie.GenerateRandomKey(32))

func main() {
    http.Handle("/", http.FileServer(http.Dir(Config.DashboardPath)))
    
    host := fmt.Sprintf("%s:%d", Config.Host, Config.Port)
    fmt.Println("Starting Jekyll Deploy on", host)
    log.Fatal(http.ListenAndServe(host, nil))
}
