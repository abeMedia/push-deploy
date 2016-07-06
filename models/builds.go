package models

import (
    "time"
    "encoding/json"
    "errors"
    t "github.com/abemedia/push-deploy/lib/db"
    "upper.io/db"
)

type build struct {
    Model
}

var Builds build

func init() {
    Builds.table = "builds"
}

type Build struct {
    ID          int `bson:"_id" json:"id"`
    ProjectID   int `json:"project_id"`
    Status      int `json:"status"`
    Start       time.Time `json:"start"`
    Finish      time.Time `json:"finish"`
    Author      string `json:"author"`
    Message     string `json:"message"`
    Logs        logs `json:"logs"`
}

type BuildList []struct {
    ID          int `bson:"_id" json:"id"`
    ProjectID   int `json:"project_id"`
    Status      int `json:"status"`
    Start       time.Time `json:"start"`
    Finish      time.Time `json:"finish"`
    Author      string `json:"author"`
    Message     string `json:"message"`
}

type logs map[string]string

// map -> string
func (l logs) MarshalDB() (interface{}, error) {
    return json.Marshal(l)
}

// string -> map
func (l *logs) UnmarshalDB(v interface{}) error {
    if _, ok := v.(string); !ok {
        return errors.New("Unsupported Value") //db.ErrUnsupportedValue
    }
    
    var d map[string]string
    json.Unmarshal([]byte(v.(string)), &d)
    *l = d
    return nil
}

func (m *build) List(id interface{}) (d BuildList, err error) {
    err = t.Table(m.table).GetList(db.Cond{"project_id":id}, &d)
    return d, err
}

func (m *build) Get(q interface{}) (d Build, err error) {
    err = t.Table(m.table).GetRow(m.query(q), &d)
    return d, err
}