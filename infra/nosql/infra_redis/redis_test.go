package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/redis/go-redis/v9"
)

var db *redis.Client

func beforeAfter(caseName string) func() {
	if db == nil {
		db = redis.NewClient(&redis.Options{
			Addr:     os_x.Getenv("redis_host") + ":" + os_x.Getenv("redis_port"),
			Username: os_x.Getenv("redis_username"),
			Password: os_x.Getenv("redis_password"),
			DB:       0,
		})
	}

	return func() {
		fmt.Println(caseName, "teardown......")
	}
}

func Test_lock_00(t *testing.T) {
	defer beforeAfter("Test_lock_00")()

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

	time.Sleep(time.Second * 5)
}
