package infra_sql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/dbs"
	_ "modernc.org/sqlite"
)

func Test_00(t *testing.T) {
	dsn := dbs.MySQLConnectWithDefaultArgs("", 3306, "root", "1121", "boot")
	db, _ := sqlx.Connect(dbs.DriverMysql, dsn)
	var tables []string
	_ = db.Select(&tables, "show tables like '%user%'")
	fmt.Println(tables)
}

func Test_lock_00(t *testing.T) {
	// dsn := dbs.MySQLConnectWithDefaultArgs("", 3306, "root", "1121", "boot")
	// db, _ := sql.Open(dbs.DriverMysql, dsn)

	home, _ := common_x.Dir()
	db, _ := dbs.SqliteConnect(home + "/data/sqlite/main.db")

	sleep := time.Second * 5
	locker := NewLocker[string](db)
	go locker.LockRun("lock_foo", time.Second*30, func() {
		fmt.Println("加锁成功，执行任务1")
		time.Sleep(sleep)
		fmt.Println("任务执行完成1")
	})
	time.Sleep(time.Second * 1)
	go locker.LockRun("lock_foo", time.Second*30, func() {
		fmt.Println("加锁成功，执行任务2")
		time.Sleep(sleep)
		fmt.Println("任务执行完成2")
	})

	time.Sleep(sleep * 2)
}
