package embedded

import (
    "database/sql"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/dbs"
    llibsql "github.com/luvx21/coding-go/coding-usage/db/sqlite/libsql"
    "github.com/tursodatabase/go-libsql"
    "path/filepath"
    "testing"
    "time"
)

func Test_00(t *testing.T) {
    dbName := "main.db"
    home, _ := common_x.Dir()
    dbPath := filepath.Join(home+"/data/sqlite/libsql", dbName)

    connector, _ := libsql.NewEmbeddedReplicaConnector(dbPath, llibsql.Url,
        libsql.WithAuthToken(llibsql.Token), libsql.WithSyncInterval(time.Second*10),
    )
    defer connector.Close()

    db := sql.OpenDB(connector)
    defer db.Close()

    rows, _ := db.Query("select * from user")
    defer rows.Close()
    dbs.PrintRows(rows)
}
