package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"
	"sort"
	"sync/atomic"
)

func main() {
	waitGroupDemo()
}

func waitGroupDemo() {
	var count atomic.Int64

	var wg conc.WaitGroup
	for i := range 10 {
		wg.Go(func() {
			if i == 7 {
				panic("异常")
			}
			count.Add(1)
		})
	}
	// recover panic
	wg.WaitAndRecover()

	fmt.Println(count.Load())
}

// poolDemo goroutine池示例
func poolDemo() {
	// 创建一个最大数量为3的goroutine池
	p := pool.New().WithMaxGoroutines(3)
	// 使用p.Go()提交5个任务
	for range 5 {
		p.Go(func() {
			fmt.Println("Q1mi")
		})
	}
	p.Wait()
}

// poolWithContextDemoCancelOnError 支持context的池
// goroutine中出错时取消context
func poolWithContextDemoCancelOnError() {
	p := pool.New().
		WithMaxGoroutines(4).
		WithContext(context.Background()).
		WithCancelOnError() // 出错时取消所有goroutine
	// 提交3个任务
	for i := range 3 {
		i := i
		p.Go(func(ctx context.Context) error {
			if i == 2 {
				return errors.New("cancel all other tasks")
			}
			<-ctx.Done()
			return nil
		})
	}
	err := p.Wait()
	fmt.Println(err)
}

// poolWithResult 执行返回结果的任务池
func poolWithResult() {
	// 创建一个任务池，其中任务返回的结果为int
	p := pool.NewWithResults[int]()
	for i := range 10 {
		i := i
		p.Go(func() int {
			return i * 2
		})
	}
	res := p.Wait()
	// 结果的顺序是不确定的, 所以这里先排序再打印
	sort.Ints(res)
	fmt.Println(res)
}

// panicDemo recover可能出现的异常
func panicDemo() {
	var pc panics.Catcher
	i := 0
	pc.Try(func() { i += 1 })
	pc.Try(func() { panic("abort!") })
	pc.Try(func() { i += 1 })

	// recover可能出现的panic
	rc := pc.Recovered()
	// 重新panic
	// pc.Repanic()

	fmt.Println(i)
	fmt.Println(rc.Value.(string))
}
