package infra_sql

import (
	"fmt"

	sql "github.com/jmoiron/sqlx"
)

const (
	DriverMysql  = "mysql"
	DriverDuckdb = "duckdb"
)

func ConnectMySQL(host string, port int, username, password, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", username, password, host, port, database)
	db, err := sql.Connect(DriverMysql, dsn)
	if err != nil {
		fmt.Printf("数据库连接失败:%v\n", err)
		return db, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, err
}

func ConnectDuckdb(dataSource string) (*sql.DB, error) {
	return sql.Open(DriverDuckdb, dataSource)
}

// func NamedExec(db *sql.DB, sql string, args ...map[string]any) (err error) {
//    for _, arg := range args {
//        _, err = db.NamedExec(sql, arg)
//    }
//    return nil
// }
//
// func NamedQuery(db *sql.DB, sql string, arg map[string]any) (err error) {
//    _, err = db.NamedQuery(sql, arg)
//    return nil
// }
