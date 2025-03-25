package pg

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/luvx21/coding-go/coding-common/dbs"

	_ "github.com/lib/pq"
)

const (
	host     = ""
	port     = 5432
	user     = ""
	password = ""
	dbname   = "postgres"
	url1     = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	url2     = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
)

var db *sql.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		_url := fmt.Sprintf(url1, host, port, user, password, dbname)
		//_url = fmt.Sprintf(url2, user, password, host, port, dbname)
		db, _ = sql.Open("postgres", _url)
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
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(url2, user, password, host, port, dbname))
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
