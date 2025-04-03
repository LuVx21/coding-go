package dbs

import (
	"database/sql"

	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

const (
	DriverDuckdb  = "duckdb"
	DriverMysql   = "mysql"
	DriverSqlite3 = "sqlite3"
	DriverSqlite  = "sqlite"
)

func PrintRows(_rows *sql.Rows) {
	columns, _ := _rows.Columns()
	values := make([]any, len(columns))
	for i := range values {
		var a any
		values[i] = &a
	}
	rows := make([][]any, 0, len(columns))
	for _rows.Next() {
		_ = _rows.Scan(values...)
		row := make([]any, len(columns))
		for i, val := range values {
			row[i] = *val.(*any)
		}
		rows = append(rows, row)
	}
	fmt_x.Println(slices_x.ToAnySliceE(columns...), rows...)
}

func ParseRows(_rows *sql.Rows) []map[string]any {
	columns, _ := _rows.Columns()
	values := make([]any, len(columns))
	for i := range values {
		var a any
		values[i] = &a
	}
	rows := make([]map[string]any, 0, len(columns))
	for _rows.Next() {
		_ = _rows.Scan(values...)
		row := make(map[string]any)
		for i, val := range values {
			row[columns[i]] = *val.(*any)
		}
		rows = append(rows, row)
	}
	return rows
}
