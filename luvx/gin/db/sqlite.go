package db

import (
    "database/sql"
    "fmt"
    "github.com/luvx21/coding-go/coding-common/common"
    _ "modernc.org/sqlite"
)

var SqliteClient *sql.DB

func init() {
    home, _ := common.Dir()
    dataSourceName := home + "/data/sqlite/main.db"
    SqliteClient, _ = GetDataSource(dataSourceName)
}

func GetDataSource(dataSourceName string) (*sql.DB, error) {
    fmt.Println("初始化Sqlite连接...", dataSourceName)
    return sql.Open("sqlite", dataSourceName)
}

func QueryForMap(db *sql.DB, sql string, args ...interface{}) ([]map[string]interface{}, error) {
    rows, err := db.Query(sql, args...)
    if err != nil {
        return nil, err
    }
    columns, _ := rows.Columns()
    columnLength := len(columns)
    cache := make([]interface{}, columnLength)
    for index, _ := range cache {
        var a interface{}
        cache[index] = &a
    }
    var list []map[string]interface{}
    for rows.Next() {
        _ = rows.Scan(cache...)

        item := make(map[string]interface{})
        for i, data := range cache {
            item[columns[i]] = *data.(*interface{})
        }
        list = append(list, item)
    }
    _ = rows.Close()
    return list, nil
}
