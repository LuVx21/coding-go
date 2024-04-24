package db

import (
    "database/sql"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/dbs"
    _ "modernc.org/sqlite"
    // _ "github.com/mattn/go-sqlite3"
)

const (
    driverName = "sqlite"
    //driverName = "sqlite3"
)

var SqliteClient *sql.DB

func init() {
    home, _ := common_x.Dir()
    dataSourceName := home + "/data/sqlite/main.db"
    SqliteClient, _ = GetDataSource(dataSourceName)
}

func GetDataSource(dataSourceName string) (*sql.DB, error) {
    defer common_x.TrackTime("初始化Sqlite连接..." + dataSourceName)()
    return sql.Open(driverName, dataSourceName)
}

func QueryForMap(db *sql.DB, sql string, args ...interface{}) ([]map[string]interface{}, error) {
    rows, err := db.Query(sql, args...)
    defer rows.Close()
    if err != nil {
        return nil, err
    }

    return dbs.ParseRows(rows), nil
}
