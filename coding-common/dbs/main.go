package dbs

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
)

var DefaultMysqlArgs = map[string]string{
	"charset":   "utf8mb4",
	"parseTime": "True",
	"loc":       "Local",
}

func MySQLConnectWithDefaultArgs(host string, port int, username, password, database string) string {
	return MySQLConnect(host, port, username, password, database, DefaultMysqlArgs)
}
func MySQLConnect(host string, port int, username, password, database string, args map[string]string) string {
	host = strings_x.FirstNonEmpty(host, "127.0.0.1")
	port = common_x.IfThen(port <= 0, 3306, port)
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	var sb strings.Builder
	sb.WriteString(s)
	if len(args) > 0 {
		sb.WriteString("?")
		values := url.Values{}
		for k, v := range args {
			values.Add(k, v)
		}
		sb.WriteString(values.Encode())
	}
	return sb.String()
}

func DuckdbConnect(dataSource string) (*sql.DB, error) {
	return sql.Open(DriverDuckdb, dataSource)
}
func SqliteConnect(dataSource string) (*sql.DB, error) {
	return sql.Open(DriverSqlite, dataSource)
}
