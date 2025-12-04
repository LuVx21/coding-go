package infra_sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

type DbLocker[T any] struct {
	Client  *sql.DB
	ownerID string // 当前实例的唯一标识

	TableCreated bool
	sqlHolder    *sqlHolder
	// mu           sync.Mutex
}

func NewLocker[T any](c *sql.DB) *DbLocker[T] {
	return NewLockerWithOwner[T](c, "default")
}
func NewLockerWithOwner[T any](c *sql.DB, ownerID string) *DbLocker[T] {
	l := &DbLocker[T]{Client: c, ownerID: ownerID}

	// 可修改为外部传参
	driverType := strings.ToLower(reflect.TypeOf(c.Driver()).String())
	if strings.Contains(driverType, "mysql") {
		l.sqlHolder = &sql_mysql
	} else if strings.Contains(driverType, "sqlite") {
		a := sql_sqlite()
		l.sqlHolder = &a
	} else if strings.Contains(driverType, "pq") {
		a := sql_postgre()
		l.sqlHolder = &a
	} else {
		slog.Error("暂不支持的数据库类型", "driverType", driverType)
		return nil
	}
	return l
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := l.Client.ExecContext(ctx, l.sqlHolder.lock_insert, key, l.ownerID, time.Now().Add(exp).UnixMilli())
	if err == nil {
		return true
	}

	var existingOwner string
	var expAt int64
	row := l.Client.QueryRowContext(ctx, l.sqlHolder.selectLock, key)
	if err := row.Scan(&existingOwner, &expAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 记录不存在(可能在检查期间被删除)，重试插入
			return l.TryLock(key, exp)
		}
		return false
	}
	// 检查锁是否已过期
	if time.Now().UnixMilli() < expAt {
		slog.Debug("现有锁未过期", "expAt", expAt)
		return false
	}

	// 锁已过期，尝试获取
	_, err = l.Client.ExecContext(ctx, l.sqlHolder.lock_update, l.ownerID, time.Now().Add(exp).UnixMilli(), key)
	if err != nil {
		slog.Error("加锁sql执行失败", "key", key, "sql", l.sqlHolder.lock_update, "error", err)
	}
	return err == nil
}

func (l *DbLocker[T]) prepare() error {
	if l.TableCreated {
		return nil
	}

	if _, err := l.Client.Exec(l.sqlHolder.ddl); err != nil {
		slog.Error("建表失败", "err", err)
		return fmt.Errorf("建表失败:%s", l.sqlHolder.ddl)
	}
	l.TableCreated = true
	return nil
}

func (l *DbLocker[T]) Unlock(key T) bool {
	// if !l.mu.TryLock() {
	// 	return false
	// }
	// defer l.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := l.Client.ExecContext(ctx, l.sqlHolder.unlock, key, l.ownerID)
	if err != nil {
		slog.Error("释放锁错误", "err", err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}
func (l *DbLocker[T]) LockRun(key T, exp time.Duration, fn func()) {
	if l.TryLock(key, exp) {
		defer l.Unlock(key)
		fn()
	} else {
		slog.Warn("加锁失败", "key", key)
	}
}

func (l *DbLocker[T]) Renew(key string, exp time.Duration) bool {
	result, err := l.Client.Exec(l.sqlHolder.renew, time.Now().Add(exp).UnixMilli(), key, l.ownerID)
	if err != nil {
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

func (l *DbLocker[T]) StartCleanup(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			_, _ = db.Exec(l.sqlHolder.cleanLock)
		}
	}()
}

type sqlHolder struct {
	ddl         string
	selectLock  string
	lock_insert string
	lock_update string
	unlock      string
	renew       string
	cleanLock   string
}

var sql_mysql = sqlHolder{
	ddl: `
		CREATE TABLE IF NOT EXISTS common_lock
		(
			id          bigint(20)   NOT NULL AUTO_INCREMENT,
			lock_key    varchar(255) NOT NULL,
			value       varchar(255) NOT NULL DEFAULT '',
			owner_id    varchar(255) NOT NULL DEFAULT '',
			expire_time bigint       NOT NULL DEFAULT 0,
			create_time timestamp    NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
			update_time timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp COMMENT '更新时间',
			PRIMARY KEY (id),
			UNIQUE KEY uk_key (lock_key)
		) ENGINE = InnoDB
		DEFAULT CHARSET = utf8mb4;
		`,
	selectLock:  "select owner_id, expire_time from common_lock where lock_key = ? limit 1 for update",
	lock_insert: "insert into common_lock (lock_key, owner_id, expire_time) values (?, ?, ?)",
	lock_update: "update common_lock set owner_id = ?, expire_time = ? where lock_key = ?",
	unlock:      "delete from common_lock where lock_key = ? and owner_id = ?",
	renew:       "update common_lock set expire_time = ? where lock_key = ? and owner_id = ?",
	cleanLock:   "delete from common_lock where expire_time < unix_timestamp() * 1000",
}

func sql_sqlite() sqlHolder {
	sql_sqlite := sql_mysql
	sql_sqlite.ddl = `
		CREATE TABLE IF NOT EXISTS common_lock
		(
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			lock_key    TEXT    NOT NULL UNIQUE,
			value       TEXT    NOT NULL DEFAULT '',
			owner_id    TEXT    NOT NULL DEFAULT '',
			expire_time INTEGER NOT NULL DEFAULT 0,
			create_time TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`
	sql_sqlite.cleanLock = "delete from common_lock where expire_time < CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER)"
	sql_sqlite.selectLock = "select owner_id, expire_time from common_lock where lock_key = ? limit 1"
	return sql_sqlite
}

func sql_postgre() sqlHolder {
	sql_postgre := sql_mysql
	sql_postgre.ddl = `
		CREATE TABLE IF NOT EXISTS common_lock (
			id          BIGSERIAL    PRIMARY KEY,
			lock_key    VARCHAR(255) NOT NULL,
			value       VARCHAR(255) NOT NULL DEFAULT '',
			owner_id    VARCHAR(255) NOT NULL DEFAULT '',
			expire_time BIGINT       NOT NULL DEFAULT 0,
			create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT uk_key UNIQUE (lock_key)
		);

		-- 为 update_time 列创建触发器以实现自动更新
		CREATE OR REPLACE FUNCTION update_modified_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.update_time = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER update_common_lock_modtime
		BEFORE UPDATE ON common_lock
		FOR EACH ROW
		EXECUTE FUNCTION update_modified_column();
	`
	sql_postgre.cleanLock = "delete from common_lock where expire_time < EXTRACT(EPOCH FROM now()) * 1000"
	return sql_postgre
}
