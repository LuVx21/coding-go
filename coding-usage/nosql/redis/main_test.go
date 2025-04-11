package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/test"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

var before = func(name string) func() {
	return test.BeforeTest(name, func() {
		if rdb == nil {
			rdb = connect()
		}
	})
}

func Test_00(t *testing.T) {
	defer before("Test_00")()
	val, _ := rdb.Get(context.TODO(), "foo").Result()
	fmt.Println("foo", "=", val)

	r := rdb.Del(context.TODO(), "lock_"+"01")
	fmt.Println(r.Result())

	b := rdb.SetNX(context.TODO(), "lock_"+"01", 1, time.Second*60)
	fmt.Println(b.Result())
}

func Test_map_00(t *testing.T) {
	defer before("Test_map_00")()
	result, _ := rdb.HGetAll(context.TODO(), "mm").Result()
	fmt.Println("mm", "=", result, result["mk"])

	v, _ := rdb.HGet(context.TODO(), "mm", "mk").Result()
	fmt.Println("mm.mk", "=", v)
}

func Test_limiter(t *testing.T) {
	defer before("Test_limiter")()

	limiter := redis_rate.NewLimiter(rdb)
	doRequest := func(i int) {
		b, _ := limiter.Allow(context.Background(), "rate_limit:user123", redis_rate.PerSecond(2))
		fmt.Printf("%v 请求 %2d: %s\n", time.Now(), i+1, common_x.IfThen(b.Allowed >= 1, "允许", "拒绝"))
	}

	// 模拟突发流量, 10个突发流量(redis服务耗时长会影响效果)
	for i := range 30 {
		doRequest(i)
		if (i+1)%10 == 0 {
			time.Sleep(2 * time.Second)
		}
	}

	time.Sleep(10 * time.Second)
}
