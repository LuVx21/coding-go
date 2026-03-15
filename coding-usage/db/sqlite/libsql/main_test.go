package libsql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/luvx21/coding-go/coding-common/dbs"
	"github.com/luvx21/coding-go/coding-usage/db"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func Test_00(t *testing.T) {
	db, _ := sql.Open("libsql", fmt.Sprintf(Url+"?authToken=%s", db.Token))
	defer db.Close()

	rows, _ := db.Query("select * from user")
	defer rows.Close()
	dbs.PrintRows(rows)
}
