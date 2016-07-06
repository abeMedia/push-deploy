// todo: make a model struct with functions overidable on a per model basis

package models

import (
    . "github.com/abemedia/push-deploy/lib/config"
    t "github.com/abemedia/push-deploy/lib/db"
    "upper.io/db"
    "encoding/json"
    "os"
    "fmt"
    "errors"
)

type projects struct {
    Model
}

var Projects projects

func init() {
    Projects.table = "projects"
}

type Project struct {
    ID      int `bson:"_id" json:"id"`
    UserID  int `json:"user_id"`
    Name    string `json:"name"`
    Repo    string `json:"repo"`
    Branch  string `json:"branch"`
    Deploy  deploy `json:"deploy"`
    Status  int `db:"-" json:"status"`
    Current Build `db:"-" json:"current"`
}

type deploy []map[string]string

// map -> string
func (u deploy) MarshalDB() (interface{}, error) {
    return json.Marshal(u)
}

// string -> map
func (u *deploy) UnmarshalDB(v interface{}) error {
    if _, ok := v.(string); !ok {
        return errors.New("Unsupported Value") //db.ErrUnsupportedValue
    }
    
    var d []map[string]string
    json.Unmarshal([]byte(v.(string)), &d)
    *u = d
    return nil
}

func (m *projects) List() (d []Project, err error) {
    err = t.Table(m.table).GetList(db.Cond{}, &d)
    if err != nil {
        return d, err
    }
    
    // avoiding joins to keep compatible with mongodb
    for i, _ := range d {
        last_build, err := Builds.Get(db.Cond{"project_id": d[i].ID})
        if err != nil {
            continue
        }
        d[i].Current = last_build
    }
    return d, nil
}

func (m *projects) Get(q interface{}) (d Project, err error) {
    err = t.Table(m.table).GetRow(m.query(q), &d)
    if err != nil {
        return d, err
    }
    last_build, _ := Builds.Get(db.Cond{"project_id": d.ID})
    if err != nil && err.Error() != "There are no more rows in this result set." {
        return d, err
    }
    d.Current = last_build
    return d, nil
}

func (m *projects) UserList(user_id int) (d []Project, err error) {
    err = t.Table(m.table).GetList(db.Cond{"user_id":user_id}, &d)
    if err != nil {
        return d, err
    }
    
    // avoiding joins to keep compatible with mongodb
    for i, _ := range d {
        last_build, err := Builds.Get(db.Cond{"project_id": d[i].ID})
        if err != nil {
            continue
        }
        d[i].Current = last_build
    }
    return d, err
}

func (m *projects) Delete(q interface{}) error {
    d, err := m.Get(q)
    if err != nil {
        return err
    }
    
    // remove all project files
    os.RemoveAll(fmt.Sprintf("%s/%d", Config.CachePath, d.ID))
    
    return t.Table(m.table).Delete(m.query(q))
}