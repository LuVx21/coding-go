package infra_sql

import (
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

type DbLocker[T any] struct {
	Client       *sql.DB
	TableCreated bool
	// mu           sync.Mutex
}

func NewLocker[T any](c *sql.DB) *DbLocker[T] {
	return &DbLocker[T]{Client: c}
}

func (l *DbLocker[T]) TryLock(key T, exp time.Duration) bool {
	// if !l.mu.TryLock() {
	// 	slog.Debug("原生锁")
	// 	return false
	// }
	// defer l.mu.Unlock()

	if err := l.prepare(); err != nil {
		slog.Debug("lock表准备", "err", err)
		return false
	}

	row := l.Client.QueryRow("select value FROM common_lock where lock_key = ? limit 1", key)
	var expAt int64
	_ = row.Scan(&expAt)
	// 现有锁未过期
	if expAt > 0 && time.Now().UnixNano() < expAt*1e6 {
		slog.Debug("现有锁未过期", "expAt", expAt*1e6)
		return false
	}

	_sql := "insert into common_lock (value, lock_key) values (?, ?)"
	if expAt > 0 {
		_sql = "update common_lock set value = ? where lock_key = ?"
	}
	endAt := time.Now().Add(exp).UnixNano() / 1e6
	_, err := l.Client.Exec(_sql, endAt, key)
	if err != nil {
		slog.Debug("加锁sql执行失败", "key", key, "sql", _sql, "error", err)
	}
	return err == nil
}

func (l *DbLocker[T]) prepare() error {
	if l.TableCreated {
		return nil
	}

	sqlText := ""
	driverType := strings.ToLower(reflect.TypeOf(l.Client.Driver()).String())
	if strings.Contains(driverType, "mysql") {
		sqlText = `
		CREATE TABLE IF NOT EXISTS common_lock
		(
			id          bigint(20)   NOT NULL AUTO_INCREMENT,
			lock_key    varchar(255) NOT NULL,
			value       bigint       NOT NULL DEFAULT 0,
			create_time timestamp    NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
			update_time timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp COMMENT '更新时间',
			PRIMARY KEY (id),
			UNIQUE KEY uk_key (lock_key)
		) ENGINE = InnoDB
		DEFAULT CHARSET = utf8mb4;
	`
	} else if strings.Contains(driverType, "sqlite") {
		sqlText = `
		CREATE TABLE IF NOT EXISTS common_lock
		(
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			lock_key    TEXT    NOT NULL UNIQUE,
			value       INTEGER NOT NULL DEFAULT 0,
			create_time TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`
	} else if strings.Contains(driverType, "pq") {
		return fmt.Errorf("暂不支持的数据库类型:%s", driverType)
	}
	if _, err := l.Client.Exec(sqlText); err != nil {
		slog.Debug("建表失败", "err", err)
		return fmt.Errorf("建表失败:%s", sqlText)
	}
	l.TableCreated = true
	return nil
}

func (l *DbLocker[T]) Unlock(key T) bool {
	// if !l.mu.TryLock() {
	// 	return false
	// }
	// defer l.mu.Unlock()

	_, err := l.Client.Exec("delete from common_lock where lock_key = ?", key)
	return err == nil
}
func (l *DbLocker[T]) LockRun(key T, exp time.Duration, fn func()) {
	if l.TryLock(key, exp) {
		defer l.Unlock(key)
		fn()
	} else {
		slog.Warn("加锁失败", "key", key)
	}
}
