package duckdb

import (
	"database/sql"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/dbs"
	_ "github.com/marcboeker/go-duckdb"
)

func Test_duckdb_c(t *testing.T) {
	home, _ := common_x.Dir()
	dataSourceName := home + "/data/duckdb/main.db"
	db, _ := sql.Open("duckdb", dataSourceName)
	defer db.Close()

	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS people (id INTEGER, name VARCHAR)`)

	_, _ = db.Exec(`INSERT INTO people VALUES (42, 'John')`)

	rows, _ := db.Query(`SELECT id, name FROM people;`)
	dbs.PrintRows(rows)
}
