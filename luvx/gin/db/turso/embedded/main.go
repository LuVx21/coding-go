package embedded

// import (
// 	"database/sql"
// 	"fmt"
// 	"path/filepath"
// 	"time"
//
// 	"github.com/luvx21/coding-go/coding-common/common_x"
// 	"github.com/tursodatabase/go-libsql"
// )
//
// func Embedded(dbname, token string) *sql.DB {
// 	home, _ := common_x.Dir()
// 	dbPath := filepath.Join(home+"/data/sqlite/libsql", dbname)
// 	url := fmt.Sprintf("libsql://%s.turso.io", dbname)
//
// 	connector, _ := libsql.NewEmbeddedReplicaConnector(dbPath, url, libsql.WithAuthToken(token),
// 		libsql.WithSyncInterval(time.Second*10),
// 	)
//
// 	return sql.OpenDB(connector)
// }
