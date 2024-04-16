package common

import (
    "fmt"
    "testing"
)

func Test_m1(t *testing.T) {
    RunCatching(task)
}

func task() {
    panic("异常")
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
