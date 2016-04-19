package db

import (
    "upper.io/db"
    "upper.io/db/mysql"
    . "github.com/abemedia/push-deploy/lib/config"
)

func init() {
    if Config.DB.Type == "mysql" {
        connect = mysql_connect
    }
}

func mysql_connect() (db.Database, error) {
    return db.Open(mysql.Adapter, &mysql.ConnectionURL{
        Address:  db.Host(Config.DB.Host),
        Database: Config.DB.Database,
        User:     Config.DB.User,
        Password: Config.DB.Password,
    })
}