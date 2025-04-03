package infra_kv

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/luvx21/coding-go/coding-common/common_x"
)

var db *badger.DB

func beforeAfter(caseName string) func() {
	if db == nil {
		home, _ := common_x.Dir()
		dbFilePath := filepath.Join(home, "data", "kv", "badger", "badger.db")

		db, _ = badger.Open(badger.DefaultOptions(dbFilePath))
	}

	return func() {
		fmt.Println(caseName, "teardown......")
	}
}

func Test_lock_00(t *testing.T) {
	defer beforeAfter("Test_lock_00")()
	defer db.Close()

	locker := NewLocker[string](db)
	go locker.LockRun("lock_foo", time.Second*30, func() {
		fmt.Println("加锁成功，执行任务1")
		time.Sleep(time.Second * 2)
		fmt.Println("任务执行完成1")
	})
	go locker.LockRun("lock_foo", time.Second*30, func() {
		fmt.Println("加锁成功，执行任务2")
		time.Sleep(time.Second * 2)
		fmt.Println("任务执行完成2")
	})

	time.Sleep(time.Second * 10)
}
