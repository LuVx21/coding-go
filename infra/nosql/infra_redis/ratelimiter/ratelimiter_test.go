package ratelimiter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/test"
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

func Test_TokenBucket_00(t *testing.T) {
	defer before("Test_TokenBucket_00")()
	fmt.Println(db.Get(context.TODO(), "foo").Val())
	testRateLimiter("令牌桶算法", TokenBucket)
}

func Test_LeakyBucket_00(t *testing.T) {
	defer before("Test_LeakyBucket_00")()
	testRateLimiter("漏桶算法", LeakyBucket)
}

func testRateLimiter(name string, algorithm AlgorithmType) {
	fmt.Printf("\n=== 测试 %s ===\n", name)

	// 创建限流器: 每秒5个请求，桶容量10
	limiter := NewRateLimiter(db, "rate_limit",
		WithAlgorithm(algorithm),
		WithRate(5),
		WithCapacity(5),
	)

	ctx := context.Background()
	identifier := "user123"

	doRequest := func(i int) {
		b, _ := limiter.Allow(ctx, identifier)

		fmt.Printf("%v 请求 %2d: %s\n", time.Now(), i+1, common_x.IfThen(b, "允许", "拒绝"))
		// time.Sleep(200 * time.Millisecond)
	}

	for i := range 30 {
		doRequest(i)
		if (i+1)%10 == 0 {
			fmt.Println("-----------------------------")
			time.Sleep(2 * time.Second)
		}
	}

	time.Sleep(10 * time.Second)
}
