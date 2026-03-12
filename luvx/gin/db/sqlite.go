package db

import (
	"database/sql"
	"log/slog"
	"time"

	"luvx/gin/common/consts"

	gorm_sqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/luvx21/coding-go/coding-common/common_x"
	_ "modernc.org/sqlite"
)

const (
	driverName = "sqlite"
	// driverName = "sqlite3"
)

var (
	SqliteClient, CookieDb *sql.DB

	FreshrssDb *gorm.DB
)

func init() {
	var err error
	SqliteClient, err = GetDataSource(consts.Home + "/data/sqlite/main.db")
	if err != nil {
		slog.Error("sqlite-SqliteClient", "err", err.Error())
	}
	go configureSQLite(SqliteClient)
	CookieDb, err = GetDataSource(consts.Home + "/data/sqlite/Cookies")
	if err != nil {
		slog.Error("sqlite-cookie", "err", err.Error())
	}

	temp, err := GetDataSource(consts.Home + "/docker/freshrss/data/users/admin/db.sqlite")
	go configureSQLite(temp)
	FreshrssDb, err = gorm.Open(gorm_sqlite.New(gorm_sqlite.Config{Conn: temp}), &gorm.Config{})
	if err != nil {
		slog.Error("sqlite-freshrss", "err", err.Error())
	}
}

func GetDataSource(dataSourceName string) (*sql.DB, error) {
	_kv, err, _ := consts.SfGroup.Do(dataSourceName, func() (any, error) {
		defer common_x.TrackTime("初始化Sqlite连接..." + dataSourceName)()
		return sql.Open(driverName, dataSourceName)
	})
	return _kv.(*sql.DB), err
}

func configureSQLite(db *sql.DB) {
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(30 * time.Minute)

	// 执行PRAGMA配置
	pragmas := []string{
		"PRAGMA journal_mode = WAL",    // 启用WAL模式，支持读写并发
		"PRAGMA synchronous = NORMAL",  // 平衡性能与安全
		"PRAGMA busy_timeout = 5000",   // 5秒锁定超时
		"PRAGMA foreign_keys = ON",     // 启用外键
		"PRAGMA cache_size = -2000",    // 2MB缓存
		"PRAGMA mmap_size = 268435456", // 256MB内存映射
		"PRAGMA wal_autocheckpoint = 800", // 当 WAL 文件达到约 1000 页时自动触发检查点
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			slog.Warn("执行失败", "命令", pragma, "错误", err)
		}
	}
}
