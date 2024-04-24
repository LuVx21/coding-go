package dbs

import (
    "context"
    "database/sql"
    "github.com/blockloop/scan"
    "github.com/luvx21/coding-go/coding-common/fmt_x"
    "github.com/luvx21/coding-go/coding-common/slices_x"
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
