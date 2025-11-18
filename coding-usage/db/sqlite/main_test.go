package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/dbs"
	//_ "modernc.org/sqlite"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const (
	//driverName = "sqlite"
	driverName = "sqlite3"
)

var db *sql.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()
		_url := home + "/data/sqlite/main.db"
		db, _ = sql.Open(driverName, _url)
	}

	return func() {
		fmt.Println(caseName, "end...")
	}
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()

	rows, _ := db.Query("SELECT * FROM user where id >= $1", 1)
	defer rows.Close()
	dbs.PrintRows(rows)
}
