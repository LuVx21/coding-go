package common_x

import (
    "fmt"
    "testing"
    "time"
)

func Test_m1(t *testing.T) {
    RunCatching(func() {
        panic("异常")
    })
    fmt.Println("后续操作1")
    r := RunCatchingReturn[string](func() string {
        panic("异常")
        return "结果"
    })
    fmt.Println(r)
    fmt.Println("后续操作2")
}

func Test_RunWithTime(t *testing.T) {
    withTime := RunWithTime("m1", func() string {
        //time.Sleep(time.Second)
        return "ok"
    })

    fmt.Println(withTime)

    time2, s := RunWithTime2("m2", func() (string, int) {
        return "ok", 1
    })
    fmt.Println(time2, s)
}

func Test_01(t *testing.T) {
    defer TrackTime1("main", time.Now())
    defer TrackTime("main1")()
    time.Sleep(time.Second * 1)
}
