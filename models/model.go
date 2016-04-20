package models

import (
    t "github.com/abemedia/push-deploy/lib/db"
    "upper.io/db"
)

type Model struct {
    table string
}

func (m *Model) query(q interface{}) db.Cond {
    switch q.(type) {
        case int: return db.Cond{"id": q.(int)} 
        case int64: return db.Cond{"id": q.(int64)} 
        case string: return db.Cond{"id": q.(string)}
        case map[string]interface{}: 
            var d db.Cond = q.(map[string]interface{})
            return d
        case nil: return db.Cond{}
    }
    return q.(db.Cond)
}

func (m *Model) List() (d []map[string]interface{}, err error) {
    err = t.Table(m.table).GetList(db.Cond{}, &d)
    return d, err
}

func (m *Model) Get(q interface{}) (d map[string]interface{}, err error) {
    err = t.Table(m.table).GetRow(m.query(q), &d)
    return d, err
}

func (m *Model) Add(d interface{}) (int64, error) {
    return t.Table(m.table).Insert(d)
}

func (m *Model) Update(q interface{}, d interface{}) error {
    return t.Table(m.table).Update(m.query(q), d)
}

func (m *Model) Delete(q interface{}) error {
    return t.Table(m.table).Delete(m.query(q))
}