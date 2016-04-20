package db

import (
    "upper.io/db"
    //"fmt"
)

var connect func() (db.Database, error)
/*
type Cond map[string]interface{}
type And []interface{}
type Or []interface{}
*/
type table struct {
    session db.Database
    table db.Collection
}

func Table(name string) *table {
    var t table
    var err error
    t.session, err = connect()
    if err != nil {
        panic(err)
    }
    
    t.table, err = t.session.Collection(name)
    if err != nil {
        panic(err)
    }
    
    return &t
}
/*
func query(query interface{}) interface{} {
    switch query.(type) {
        case db.Cond: 
            var q map[string]interface{}
            var result db.Cond
            q = query.()
            result = q
            return result
        case db.Or: 
            var q []interface{}
            var result db.Or
            q = query
            result = q
            return result
        case db.And: 
            var q []interface{}
            var result db.And
            q = query
            result = q
            return result
    }
}
*/
func (d *table) GetList(q interface{}, result interface{}) error {
    defer d.session.Close()
    return d.table.Find(q).All(result)
}

func (d *table) GetRow(q interface{}, result interface{}) error {
    defer d.session.Close()
    return d.table.Find(q).Limit(1).Sort("-id").One(result)
}

func (d *table) Insert(data interface{}) (int64, error) {
    defer d.session.Close()
    id, err := d.table.Append(data)
    if err != nil {
        return 0, err
    }
    
    return id.(int64), nil
}

func (d *table) Update(q interface{}, data interface{}) error {
    defer d.session.Close()
    return d.table.Find(q).Update(data)
}

func (d *table) Delete(q interface{}) error {
    defer d.session.Close()
    return d.table.Find(q).Remove()
}