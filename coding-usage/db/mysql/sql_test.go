package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/dbs"
	"github.com/luvx21/coding-go/infra/infra_sql"
)

var db *sql.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		db, _ = sql.Open(dbs.DriverMysql, dbs.MySQLConnectWithDefaultArgs("", 53307, "root", "1121", "boot"))
	}
	return func() {
		fmt.Println(caseName, "teardown......")
	}
}

func Test_ddl(t *testing.T) {
	defer beforeAfter("Test_ddl")()
	sqlText := `
        CREATE TABLE IF NOT EXISTS test (
          id bigint(20) NOT NULL AUTO_INCREMENT,
          user_name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
          password varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
          age int(11) DEFAULT NULL,
          PRIMARY KEY (id)
        ) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8;
    `
	_, err := db.Exec(sqlText)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_insert(t *testing.T) {
	defer beforeAfter("Test_insert")()

	sqlText := `INSERT INTO test(user_name, password, age) VALUES (?, ?, ?);`
	rs, _ := db.Exec(sqlText, "foo3", "bar", 3)
	rowCount, _ := rs.RowsAffected()
	fmt.Printf("inserted %d rows\n", rowCount)
}

func Test_select(t *testing.T) {
	defer beforeAfter("Test_select")()

	rowsMap, _ := infra_sql.RowsMap(context.Background(), db, "SELECT * FROM test where id = ? limit 10", 50)
	for _, m := range rowsMap {
		for k, v := range m {
			fmt.Println(k, "=", cast_x.ToString(v))
		}
	}
}
