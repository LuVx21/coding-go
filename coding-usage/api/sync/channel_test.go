package main

import (
    "fmt"
    "log"
    "testing"
    "time"
)

func Producer(channel chan<- int) {
    for i := 0; i < 10; i++ {
        channel <- i //写入
        fmt.Println("发送:", i)
        time.Sleep(time.Second)
    }
}

func Consumer(channel <-chan int) {
    for i := 0; i < 10; i++ {
        v := <-channel // 读出
        fmt.Println("接收:", v)
        time.Sleep(time.Second)
    }
}

func Test_channel(t *testing.T) {
    channel := make(chan int, 88)
    go Producer(channel)
    go Consumer(channel)
    time.Sleep(12 * time.Second)
}

func TestTask1(t *testing.T) {
    ch := make(chan struct{}) // 初始化 chan
    log.Println("point0")
    go func() {
        log.Println("point---")
        Task()
        ch <- struct{}{} // 发送到 chan
    }()
    log.Println("point1")
    _ = <-ch // 从 chan 获取,阻塞
    log.Println("point2")
    log.Println("main done")

    time.Sleep(time.Second * 4)
}

func Test_c_01(t *testing.T) {
    r1 := make(chan string, 1)
    r1 <- "a"
    <-r1
    r1 <- "b"
    close(r1)
    for s := range r1 {
        fmt.Println(s)
    }
}
