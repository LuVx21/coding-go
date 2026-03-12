package turso

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/luvx21/coding-go/coding-usage/db"
	turso "turso.tech/database/tursogo"
)

var (
	path = os.ExpandEnv("${HOME}/data/sqlite/turso/sqlite.db")
)

const (
	ddl_create = "CREATE TABLE IF NOT EXISTS go_turso (foo INTEGER, bar TEXT)"
	dml_insert = "INSERT INTO go_turso (foo, bar) VALUES (?, ?)"
	dml_select = "SELECT * FROM go_turso"
)

func Test_turso_00(t *testing.T) {
	conn, _ := sql.Open("turso", path)
	defer conn.Close()

	_, _ = conn.Exec(ddl_create)

	stmt, _ := conn.Prepare(dml_insert)
	defer stmt.Close()
	_, _ = stmt.Exec(42, "turso")

	rows, _ := conn.Query(dml_select)
	defer rows.Close()
	for rows.Next() {
		var a int
		var b string
		_ = rows.Scan(&a, &b)
		fmt.Printf("%d, %s\n", a, b) // 42, turso
	}
}

func Test_turso_01(t *testing.T) {
	ctx := context.Background()

	db, err := turso.NewTursoSyncDb(ctx, turso.TursoSyncDbConfig{
		Path:      path,
		RemoteUrl: fmt.Sprintf("https://%s.%s.turso.io", db.Db, db.Region),
		AuthToken: db.Token,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	conn, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, _ = conn.ExecContext(ctx, ddl_create)

	stmt, err := conn.PrepareContext(ctx, dml_insert)
	t.Log("stmt", err)
	defer stmt.Close()
	_, _ = stmt.ExecContext(ctx, 43, "turso-remote")

	// Push local commits to remote
	_ = db.Push(ctx)

	// Pull new changes from remote into local
	_, _ = db.Pull(ctx)

	rows, _ := conn.QueryContext(ctx, dml_select)
	defer rows.Close()
	for rows.Next() {
		var a int
		var b string
		_ = rows.Scan(&a, &b)
		fmt.Printf("%d, %s\n", a, b) // 42, turso
	}

	// Optional: inspect and manage sync state
	stats, err := db.Stats(ctx)
	if err != nil {
		log.Println("Stats unavailable:", err)
	} else {
		log.Println("Current revision:", stats.NetworkReceivedBytes)
	}

	_ = db.Checkpoint(ctx)
}
