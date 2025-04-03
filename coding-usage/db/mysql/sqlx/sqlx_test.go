package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/dbs"
	cur "github.com/luvx21/coding-go/coding-usage/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type User struct {
	Id         uint32
	UserName   string         `db:"user_name"`
	Password   sql.NullString `db:"password"`
	Age        int8
	Birthday   sql.NullTime
	UpdateTime time.Time `db:"update_time"`
}

func beforeAfter(caseName string) func() {
	if db == nil {
		db = sqlx.MustConnect(dbs.DriverMysql, dbs.MySQLConnectWithDefaultArgs(cur.MysqlHost, cur.MysqlPort, cur.MysqlUser, cur.MysqlPassword, cur.MysqlDbname))
	}
	return func() {
		fmt.Println(caseName, "teardown......")
	}
}

func Test_sqlx_00(t *testing.T) {
	defer beforeAfter("Test_sqlx_00")()

	var user User
	err := db.Get(&user, "SELECT * FROM user WHERE id = 1")
	fmt.Printf("%#v %v\n", err, user)

	var users []User
	_ = db.Select(&users, "SELECT * FROM user limit 10")
	fmt.Println(users[0])

	rows, _ := db.NamedQuery(`SELECT * FROM user WHERE user_name=:aaa order by id`, map[string]any{"aaa": "foo"})
	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v\n", user)
	}
}
