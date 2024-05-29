package main

import (
    "database/sql"
    "fmt"
    "github.com/luvx21/coding-go/coding-usage/db"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

type User struct {
    Id         uint32
    UserName   string         `db:"user_name"`
    Password   sql.NullString `db:"password"`
    Age        int8
    Birthday   time.Time
    UpdateTime time.Time
}

var (
    host     = "luvx"
    port     = db.MysqlPort
    user     = "root"
    password = db.MysqlPassword
    dbname   = "boot"
    url      = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
    db := sqlx.MustConnect("mysql", fmt.Sprintf(url, user, password, host, port, dbname))

    var users []User
    _ = db.Select(&users, "SELECT * FROM user limit 10")
    fmt.Println(users)
    //fmt.Printf("%#v\n", users[0])

    user := User{}
    _ = db.Get(&user, "SELECT * FROM user WHERE id = 1")
    fmt.Printf("%#v\n", user)

    rows, _ := db.NamedQuery(`SELECT * FROM user WHERE user_name=:aaa`, map[string]any{"aaa": "foo"})
    for rows.Next() {
        err := rows.StructScan(&user)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Printf("%#v\n", user)
    }
}
