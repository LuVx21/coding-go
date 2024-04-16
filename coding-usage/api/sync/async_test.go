package main

import (
    "fmt"
    "log"
    "testing"
    "time"
)

func Test_block(t *testing.T) {
    a := 1
    b := 2

    // 从管道中接收结果，这一步是阻塞的，因为在等待结果的产出
    sum := <-addAsync(a, b)
    fmt.Println("time main", time.Now())
    log.Println(sum)
}

func addAsync(a int, b int) chan int {
    // 使用管道接收结果，注意需要设置一个缓冲位，否则没有取结果的话这个 goroutine 会被阻塞
    resultChan := make(chan int, 1)
    go func() {
        fmt.Println("time 2", time.Now())
        // 在新的 goroutine 中计算结果，并将结果发送到管道
        time.Sleep(time.Second)
        resultChan <- a + b
    }()
    fmt.Println("time 1", time.Now())
    return resultChan
}

func Test_unblock(t *testing.T) {
    a := 1
    b := 2

    // 调用方法的时候加上回调函数
    // 这个回调函数会在得到结果之后执行
    addWithCallback(a, b, func(sum int) {
        fmt.Println("time 3", time.Now())
        log.Println(sum)
    })

    fmt.Println("time main", time.Now())
    // 防止 main goroutine 比异步任务的 goroutine 先退出
    time.Sleep(2 * time.Second)
}

func addWithCallback(a int, b int, callback func(sum int)) {
    go func() {
        // 在新的 goroutine 中计算结果，并将结果传递给回调函数
        time.Sleep(time.Second)
        sum := a + b
        fmt.Println("time 2", time.Now())
        callback(sum)
    }()
}
