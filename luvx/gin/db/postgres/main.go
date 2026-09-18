package postgres

import (
	"fmt"
	"log/slog"
	"luvx/gin/config"
	"os"
	"strings"
	"sync"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	PostgreCli         = sync.OnceValue(func() *gorm.DB { return createPostgreSQLCli("") })
	PostgreCliFreshRss = sync.OnceValue(func() *gorm.DB { return createPostgreSQLCli("freshrss") })
)

func createPostgreSQLCli(dbname string) *gorm.DB {
	defer common_x.TrackTime("初始化PostgreSQL连接...")()

	c := config.AppConfig.PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Shanghai", c.Username, c.Password, c.Host, c.Port, common_x.IfThen(dbname == "", c.Dbname, dbname))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("连接数据库失败", "error", err)
		os.Exit(1)
	}
	return db
}

func ToPgJsonPath(parts []string) string {
	escaped := make([]string, len(parts))
	for i, p := range parts {
		// 双引号包裹，内部的双引号和反斜杠转义
		p = strings.ReplaceAll(p, `\`, `\\`)
		p = strings.ReplaceAll(p, `"`, `\"`)
		escaped[i] = `"` + p + `"`
	}
	return "{" + strings.Join(escaped, ",") + "}"
}
