package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/test"
	"github.com/luvx21/coding-go/infra/nosql/infra_redis/frequencylimiter"
	"github.com/redis/go-redis/v9"
)

var db *redis.Client

var before = func(name string) func() {
	return test.BeforeTest(name, func() {
		if db == nil {
			db = redis.NewClient(&redis.Options{
				Addr:     os_x.Getenv("redis_host") + ":" + os_x.Getenv("redis_port"),
				Username: os_x.Getenv("redis_username"),
				Password: os_x.Getenv("redis_password"),
				DB:       0,
			})
		}
	})
}

func Test_lock_00(t *testing.T) {
	defer before("Test_lock_00")()

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

func Test_frequencylimiter_00(t *testing.T) {
	defer before("Test_frequencylimiter_00")()

	limiter := frequencylimiter.NewFrequencyLimiter(db, "frequency_limit")

	ctx := context.Background()
	identifier := "user123"

	five_time_pre_day := func(k string, i int32) (bool, error) { return limiter.DecrTimesInDay(ctx, k, i, 5) }
	doRequest := func() {
		b, _ := five_time_pre_day(identifier, 1)

		v := db.Get(t.Context(), "frequency_limit"+":"+identifier).Val()
		fmt.Println("结果:", b, v)
	}

	doRequest()

	time.Sleep(3 * time.Second)
}
