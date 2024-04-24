package pg

import (
    "database/sql"
    "fmt"
    "github.com/luvx21/coding-go/coding-common/dbs"
    "testing"

    _ "github.com/lib/pq"
)

const (
    host     = ""
    port     = 5432
    user     = ""
    password = ""
    dbname   = "postgres"
    url1     = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
    url2     = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
)

var db *sql.DB

func beforeAfter(caseName string) func() {
    if db == nil {
        _url := fmt.Sprintf(url1, host, port, user, password, dbname)
        //_url = fmt.Sprintf(url2, user, password, host, port, dbname)
        db, _ = sql.Open("postgres", _url)
    }

    return func() {
        fmt.Println(caseName, "end...")
    }
}

func Test_00(t *testing.T) {
    defer beforeAfter("Test_00")()

    rows, _ := db.Query("SELECT * FROM boot.t_user where id >= $1", 1)
    defer rows.Close()
    dbs.PrintRows(rows)
}
