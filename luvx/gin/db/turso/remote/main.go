package remote

import (
	"database/sql"
	"fmt"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func Remote(dbname, token string) *sql.DB {
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbname, token)
	db, err := sql.Open("libsql", url)
	if err != nil {
		panic(err)
	}
	return db
}
