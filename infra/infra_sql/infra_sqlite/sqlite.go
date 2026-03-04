package infra_sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const (
	Drive, Drive3 = "sqlite", "sqlite3"
)

func Load(path string) (*sql.DB, error) {
	return sql.Open(Drive, path)
}
