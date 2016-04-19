package models

import (
    t "github.com/abemedia/push-deploy/lib/db"
    "upper.io/db"
)

type queue struct {
    Model
}

var Queues queue

func init() {
    Queues.table = "queue"
}

type Queue struct {
    ID          int `bson:"_id" json:"id"`
    ProjectID   int `json:"project_id"`
}

func (m *queue) List() (d []Queue, err error) {
    err = t.Table(m.table).GetList(db.Cond{}, &d)
    return d, err
}

func (m *queue) Get(q interface{}) (d Queue, err error) {
    err = t.Table(m.table).GetRow(m.query(q), &d)
    return d, err
}