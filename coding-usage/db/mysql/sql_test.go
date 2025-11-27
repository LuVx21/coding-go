package mysql

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"testing"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/dbs"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/infra/infra_sql"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db     *sql.DB
	gormDB *gorm.DB
)

func beforeAfter(caseName string) func() {
	if db == nil {
		db, _ = sql.Open(dbs.DriverMysql, dbs.MySQLConnectWithDefaultArgs("", 53306, "root", "1121", "boot"))
		// db = tidb()
	}

	if gormDB == nil {
		gormDB, _ = gorm.Open(
			gorm_mysql.New(gorm_mysql.Config{Conn: db}),
			&gorm.Config{SkipDefaultTransaction: true}, // 建议关闭默认事务
		)
	}

	return func() {
		fmt.Println(caseName, "teardown......")
	}
}

func tidb() *sql.DB {
	tidbHost := os_x.Getenv("TIDB_HOST")
	tidbUsername := os_x.Getenv("TIDB_USERNAME")
	tidbPassword := os_x.Getenv("TIDB_PASSWORD")
	tidbPort := os_x.Getenv("TIDB_PORT")

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: tidbHost,
	})

	dsn := dbs.MySQLConnect(tidbHost, cast_x.ToInt(tidbPort), tidbUsername, tidbPassword, "boot", map[string]string{"tls": "tidb"})
	db, err := sql.Open(dbs.DriverMysql, dsn)
	if err != nil {
		slog.Error("错误", "err", err)
		return nil
	}
	return db
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

type Cookie struct {
	ID      uint   `gorm:"primaryKey"`
	HostKey string `gorm:"size:255;not null"`
	Name    string `gorm:"size:255"`
	Value   string `gorm:"size:255"`
}

func (Cookie) TableName() string {
	return "cookie"
}

func Test_create_table(t *testing.T) {
	defer beforeAfter("Test_select")()

	gormDB, err := gorm.Open(
		gorm_mysql.New(gorm_mysql.Config{Conn: db}),
		&gorm.Config{SkipDefaultTransaction: true}, // 建议关闭默认事务
	)
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移（创建表）
	err = gormDB.AutoMigrate(&Cookie{})
	if err != nil {
		panic(fmt.Sprintf("AutoMigrate failed: %v", err))
	}

	fmt.Println("Table created successfully!")
}
