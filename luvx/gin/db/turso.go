package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"

	"luvx/gin/common/consts"

	"github.com/luvx21/coding-go/coding-common/common_x"
	turso "turso.tech/database/tursogo"

	"luvx/gin/config"
)

var (
	Turso = sync.OnceValue(createTursoCli)
)

func createTursoCli() *sql.DB {
	defer common_x.TrackTime("初始化Turso连接...")()

	c := config.AppConfig.Turso
	// Turso = embedded.Embedded(c.Dbname, c.Token)
	// Turso = remote.Remote(c.Dbname, c.Token)

	db, err := turso.NewTursoSyncDb(context.TODO(), turso.TursoSyncDbConfig{
		Path:      filepath.Join(consts.Home+"/data/sqlite/turso", c.Dbname, "sqlite.db"),
		RemoteUrl: fmt.Sprintf("https://%s.%s.turso.io", c.Dbname, "aws-ap-northeast-1"),
		AuthToken: c.Token,
	})
	if err != nil {
		slog.Error("turso", "err", err.Error())
		panic(err)
	}

	cli, err := db.Connect(context.TODO())
	if err != nil {
		slog.Error("turso db", "err", err.Error())
		panic(err)
	}
	return cli
}
