package api

import (
    "fmt"
    "log"
    "testing"
)

func TestPanic(t *testing.T) {
    // 没有这段,panic会异常退出
    defer func() {
        if r := recover(); r != nil {
            t.Log("结果:", r)
        }
    }()

    // panic只对当前goroutine的defer有效!
    //go func() {
    panic("异常了")
    //}()
    log.Print("end")
}

func Test_panic_00(t *testing.T) {
    f := func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("异常")
            }
        }()
        panic("异常")
        fmt.Println("执行")
    }
    f()
    fmt.Println("后续操作")
}
