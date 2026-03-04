package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	retry_1 "github.com/sethvargo/go-retry"
)

// 生成一个随机数，如果小于50，则返回错误，如果大于50，则返回这个数
// 可调节数字控制异常出现的概率
func retryable() (num int, err error) {
	num = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
	fmt.Println(times_x.TimeNowDateSecond(), "生成随机数:", num)
	if num < 70 {
		return 0, fmt.Errorf("人造错误")
	}
	return
}

func retryable_no_data() error {
	_, err := retryable()
	return err
}

func Test_retry_00(t *testing.T) {
	// 使用重试策略进行重试
	err := retry.New(retry.Delay(1*time.Second), retry.Attempts(5), retry.LastErrorOnly(true)).Do(retryable_no_data)
	fmt.Println(err)

	num, err := retry.NewWithData[int](retry.Delay(1*time.Second), retry.Attempts(5), retry.LastErrorOnly(true)).Do(retryable)

	if err != nil {
		fmt.Println("5次重试后仍错误")
	} else {
		fmt.Println(num)
	}
}

func Test_retry_go_retry_00(t *testing.T) {
	f1 := func(ctx context.Context) (int, error) {
		n, err := retryable()
		// 仅retryableError触发重试
		return n, common_x.IfThen(err != nil, retry_1.RetryableError(err), nil)
	}
	f := func(ctx context.Context) error {
		_, e := f1(ctx)
		return e
	}

	err := retry_1.Do(context.Background(), retry_1.WithMaxRetries(5, retry_1.NewFibonacci(1*time.Second)), f)
	fmt.Println("错误:", err)

	// backoff 不可复用
	data, err := retry_1.DoValue(context.Background(), retry_1.WithMaxRetries(5, retry_1.NewFibonacci(1*time.Second)), f1)
	fmt.Println(data, "错误:", err)
}
