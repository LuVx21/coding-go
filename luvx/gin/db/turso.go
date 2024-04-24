package db

import (
    "database/sql"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "luvx/gin/config"
    "luvx/gin/db/turso/remote"
)

var Turso *sql.DB

func init() {
    defer common_x.TrackTime("初始化Turso连接...")()

    c := config.AppConfig.Turso
    //Turso = embedded.Embedded(c.Dbname, c.Token)
    Turso = remote.Remote(c.Dbname, c.Token)
}
