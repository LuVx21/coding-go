package main

import (
    "fmt"
    "github.com/sourcegraph/conc/iter"
    "github.com/sourcegraph/conc/stream"
    "testing"
    "time"
)

func Test_iterator(t *testing.T) {
    input := []int{1, 2, 3, 4}
    // 创建一个最大goroutine个数为输入元素一半的迭代器
    iterator := iter.Iterator[int]{
        MaxGoroutines: len(input) / 2,
    }

    iterator.ForEach(input, func(v *int) {
        if *v%2 != 0 {
            *v = -1
        }
    })

    fmt.Println(input)
}

func Test_mapper(t *testing.T) {
    input := []int{1, 2, 3, 4}
    // 创建一个最大goroutine个数为输入元素一半的映射器
    mapper := iter.Mapper[int, bool]{
        MaxGoroutines: len(input) / 2,
    }

    results := mapper.Map(input, func(v *int) bool { return *v%2 == 0 })
    fmt.Println(results)
}

func Test_stream(t *testing.T) {
    times := []int{20, 52, 16, 45, 4, 80}

    s := stream.New()
    for _, millis := range times {
        dur := time.Duration(millis) * time.Millisecond
        // 提交任务
        s.Go(func() stream.Callback {
            time.Sleep(dur)
            // 虽然上一行通过sleep增加了时间
            // 但最终结果仍按任务提交（s.Go）的顺序打印
            return func() { fmt.Println(dur) }
        })
    }
    s.Wait()
}
