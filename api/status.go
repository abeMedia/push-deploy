package api

import (
    "fmt"
    t "github.com/abemedia/push-deploy/lib/db"
    "upper.io/db"
    "net/http"
    "time"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func ProjectStatus(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Print("upgrade:", err)
        return
    }
    
    defer c.Close()
    
    id := mux.Vars(r)["id"]
    
    // todo: make status come from channel in agent
    var p map[string]interface{}
    err = t.Table("builds").GetRow(db.Cond{"project_id": id}, &p)
    
    status, _ := p["status"].(string)
    
    for {
        t.Table("builds").GetRow(db.Cond{"project_id": id}, &p)
        if d, ok := p["status"].(string); ok && d != status {
            status = d
            err = c.WriteMessage(websocket.TextMessage, []byte(status))
            if err != nil {
                break
            }
        } 
        /*
        // break if socket is closed by client
        _, _, err = c.ReadMessage()
        if err != nil {
            fmt.Print("upgrade:", err)
            break
        }
        */
        
        time.Sleep(1 * time.Second)
    }
  
    /*
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Fprintln(w, "hello")
    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    
    fmt.Println(id)
    channel := agent.GetStatusChannel(id)
    fmt.Println(channel)
    for {
        select {
            case status := <-channel : 
                if err = conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(status))); err != nil {
                    fmt.Println(err)
                }
                fmt.Printf("%s", r)
        }
    }
    */
}

func ProjectsStatus(w http.ResponseWriter, r *http.Request) {
    // upgrade to websocket connection
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Print("upgrade:", err)
        return
    }
    defer c.Close()
    
    user_id := mux.Vars(r)["user_id"]
    
    // todo: make status come from channel in agent
    
    // get user's project list (not using joins for mongodb compatibility)
    var p []map[string]interface{}
    err = t.Table("projects").GetList(db.Cond{"user_id": user_id}, &p)
    if err != nil {
        fmt.Println(err)
        return
    }
    var cond db.Or
    for i, _ := range p {
        cond = append(cond, db.Cond{"project_id": p[i]["id"].(string)})
    }
    
    err = t.Table("builds").GetList(cond, &p)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    status := make(map[int]string)
    for i, _ := range p {
        if s, ok := p[i]["status"].(string); ok {
            status[i] = s
        }
    }
    
    for {
        t.Table("builds").GetList(cond, &p)
        for i, _ := range p {
            if d, ok := p[i]["status"].(string); ok && d != status[i] {
                status[i] = d
                err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d:%s", p[i]["id"], d)))
                if err != nil {
                    fmt.Println("write:", err)
                    break
                }
            }
        }
        
        if err != nil {
            break
        }
        
        /*
        // break if socket is closed by client
        _, _, err = c.ReadMessage()
        if err != nil {
            fmt.Print("upgrade:", err)
            break
        }
        */
        
        time.Sleep(5 * time.Second)
    }
}

