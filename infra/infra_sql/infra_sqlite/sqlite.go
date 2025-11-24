package infra_sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const (
	drive = "sqlite"
)

func Load(path string) (*sql.DB, error) {
	return sql.Open(drive, path)
}
