package infra_sql

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    sql "github.com/jmoiron/sqlx"
)

func Connect(host string, port int, username, password, database string) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", username, password, host, port, database)
    db, err := sql.Connect("mysql", dsn)
    if err != nil {
        fmt.Printf("数据库连接失败:%v\n", err)
        return db, err
    }

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(10)
    return db, err
}

//func NamedExec(db *sql.DB, sql string, args ...map[string]any) (err error) {
//    for _, arg := range args {
//        _, err = db.NamedExec(sql, arg)
//    }
//    return nil
//}
//
//func NamedQuery(db *sql.DB, sql string, arg map[string]any) (err error) {
//    _, err = db.NamedQuery(sql, arg)
//    return nil
//}
