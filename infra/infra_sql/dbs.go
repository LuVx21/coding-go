package infra_sql

import (
	"context"
	"database/sql"

	"github.com/blockloop/scan"
)

func Rows[T any](ctx context.Context, db *sql.DB, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var result []T
	err = scan.Rows(&result, rows)

	return result, err
}

func Row[T any](ctx context.Context, db *sql.DB, query string, args ...any) (T, error) {
	var result T

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return result, err
	}

	err = scan.Row(&result, rows)

	return result, err
}

func RowsMap(ctx context.Context, db *sql.DB, query string, args ...any) ([]map[string]any, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]any, len(colNames))
	colPtrs := make([]any, len(colNames))

	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}

	var ret []map[string]any
	for rows.Next() {
		err = rows.Scan(colPtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range cols {
			row[colNames[i]] = col
		}
		ret = append(ret, row)
	}

	return ret, nil
}

func RowMap(ctx context.Context, db *sql.DB, query string, args ...any) (map[string]any, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]any, len(colNames))
	colPtrs := make([]any, len(colNames))

	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(colPtrs...)
	if err != nil {
		return nil, err
	}

	row := make(map[string]any)
	for i, col := range cols {
		row[colNames[i]] = col
	}

	return row, nil
}
