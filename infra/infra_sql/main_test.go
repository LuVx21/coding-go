package infra_sql

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_00(t *testing.T) {
	db, _ := ConnectMySQL("127.0.0.1", 3306, "root", "1121", "boot")
	var tables []string
	_ = db.Select(&tables, "show tables like '%user%'")
	fmt.Println(tables)
}
