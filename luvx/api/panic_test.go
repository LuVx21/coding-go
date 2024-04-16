package api

import (
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
