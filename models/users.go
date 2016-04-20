// todo: make a model struct with functions overidable on a per model basis

package models

import (
    "fmt"
    t "github.com/abemedia/push-deploy/lib/db"
    "upper.io/db"
    "golang.org/x/crypto/bcrypt"
)

type users struct {
    Model
}

var Users users

func init() {
    Users.table = "users"
}

type User struct {
    ID          int `bson:"_id" json:"id"`
    Name        string `json:"name"`
    Email       string `json:"email"`
    Username    string `json:"username"`
    Password    []byte `db:"password" json:"-"`
    Admin       bool `json:"admin"`
}

func (m *users) Auth(username, password string) (User, error) {
    var u User
    err := t.Table(m.table).GetRow(db.Cond{"username":username}, &u)
    if err != nil {
        return u, err
    }
    err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
    if err != nil {
        return User{}, err
    }
    return u, nil
}

func (m *users) Add(u User) (int64, error) {
    // Hashing the password with the default cost of 10
    hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return 0, err
    }
    u.Password = hash
    return t.Table(m.table).Insert(u)
}

func (m *users) Update(q interface{}, u User) error {
    if u.Password != nil {
        // Hashing the password with the default cost of 10
        hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = hash
    }
    return t.Table(m.table).Update(m.query(q), u)
}

func (m *users) List() ([]User, error) {
    var u []User
    err := t.Table(m.table).GetList(db.Cond{}, &u)
    if err != nil {
        return []User{}, err
    }
    return u, err
}

func (m *users) Get(q interface{}) (User, error) {
    var u User
    err := t.Table(m.table).GetRow(m.query(q), &u)
    if err != nil {
        return User{}, err
    }
    return u, nil
}

func (m *users) Delete(q interface{}) error {
    var u User
    err := t.Table(m.table).GetRow(m.query(q), &u)
    if err != nil {
        return err
    }
    err = Projects.Delete(db.Cond{"user_id":u.ID})
    if err != nil && err.Error() != "There are no more rows in this result set." {
        fmt.Printf("%s",err)
        return err
    }
    return t.Table(m.table).Delete(db.Cond{"id":u.ID})
}