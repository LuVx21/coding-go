package pg

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/luvx21/coding-go/coding-common/dbs"
	"github.com/luvx21/coding-go/coding-common/jsons"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = ""
	port     = 5432
	user     = ""
	password = ""
	dbname   = "postgres"

	url = "postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Shanghai"
)

var (
	db  *sql.DB
	gdb *gorm.DB
)

func beforeAfter(caseName string) func() {
	if db == nil {
		_url := fmt.Sprintf(url, user, password, host, port, dbname)
		db, _ = sql.Open("pgx", _url)
		gdb, _ = gorm.Open(postgres.Open(_url), &gorm.Config{})
	}

	return func() {
		fmt.Println(caseName, "end...")
	}
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()

	rows, _ := db.Query("SELECT * FROM boot.t_user where id >= $1", 1)
	defer rows.Close()
	dbs.PrintRows(rows)
}

func Test_01(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(url, user, password, host, port, dbname))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
}
