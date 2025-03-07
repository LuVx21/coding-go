package duckdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x"
	_ "github.com/marcboeker/go-duckdb"
)

func Test_duckdb_c(t *testing.T) {
	home, _ := common_x.Dir()
	dataSourceName := home + "/data/sqlite/identifier.db"
	db, _ := sql.Open("duckdb", dataSourceName)
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE people (id INTEGER, name VARCHAR)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO people VALUES (42, 'John')`)
	if err != nil {
		log.Fatal(err)
	}

	var (
		id   int
		name string
	)
	row := db.QueryRow(`SELECT id, name FROM people`)
	err = row.Scan(&id, &name)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no rows")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("id: %d, name: %s\n", id, name)
}
