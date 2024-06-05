package db

import (
    "database/sql"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "luvx/gin/common/consts"
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
    _kv, err, _ := consts.SfGroup.Do(dataSourceName, func() (any, error) {
        defer common_x.TrackTime("初始化Sqlite连接..." + dataSourceName)()
        return sql.Open(driverName, dataSourceName)
    })
    return _kv.(*sql.DB), err
}
