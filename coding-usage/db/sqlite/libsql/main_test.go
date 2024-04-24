package libsql

import (
    "database/sql"
    "fmt"
    "github.com/luvx21/coding-go/coding-common/dbs"
    _ "github.com/tursodatabase/libsql-client-go/libsql"
    "testing"
)

func Test_00(t *testing.T) {
    db, _ := sql.Open("libsql", fmt.Sprintf(Url+"?authToken=%s", Token))
    defer db.Close()

    rows, _ := db.Query("select * from user")
    defer rows.Close()
    dbs.PrintRows(rows)
}
