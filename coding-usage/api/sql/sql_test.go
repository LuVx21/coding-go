package sql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:1121@tcp(luvx:53306)/boot?charset=utf8&parseTime=true")
	db.SetMaxOpenConns(2_000)
	db.SetMaxIdleConns(1_000)
	db.Ping()
}

func Test_ddl(t *testing.T) {
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
	sqlText := `INSERT INTO test(user_name, password, age) VALUES ('foo3', 'bar', 3);`
	rs, _ := db.Exec(sqlText)
	rowCount, _ := rs.RowsAffected()
	fmt.Printf("inserted %d rows\n", rowCount)
}

func Test_select(t *testing.T) {
	rows, _ := db.Query("SELECT * FROM test limit 10")
	defer rows.Close()

	columns, _ := rows.Columns()

	scanArgs := make([]any, len(columns))
	values := make([]any, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]any)
	t.Log(columns)
	for rows.Next() {
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		t.Log(record)
	}
}
