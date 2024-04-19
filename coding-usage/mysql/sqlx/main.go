package main

import (
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

type User struct {
    Id         uint
    UserName   string
    Password   string
    Age        int8
    Birthday   time.Time
    UpdateTime time.Time
}

func main() {
    db := sqlx.MustConnect("mysql", "root:xxx@tcp(xxx:3306)/boot?charset=utf8mb4&parseTime=True&loc=Local")

    var users []User
    _ = db.Select(&users, "SELECT * FROM user limit 10")
    fmt.Println(users)
    //fmt.Printf("%#v\n", users[0])

    user := User{}
    _ = db.Get(&user, "SELECT * FROM user WHERE id = 1")
    fmt.Printf("%#v\n", user)

    rows, _ := db.NamedQuery(`SELECT * FROM user WHERE user_name=:aaa`, map[string]interface{}{"aaa": "foo"})
    for rows.Next() {
        err := rows.StructScan(&user)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Printf("%#v\n", user)
    }
}
