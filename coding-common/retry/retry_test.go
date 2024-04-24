package retry

import (
    "fmt"
    "testing"
    "time"
)

func mayFail() {
    fmt.Println("执行操作...")
    if time.Now().UnixNano()%3 == 0 {
        panic("发生异常")
    }
    fmt.Println("执行结束...")
}

func Test_01(t *testing.T) {
    f := func() string {
        mayFail()
        return "结果"
    }
    retry, err := SupplyWithRetry("", f, 5, time.Second*3)
    if err != nil {
        fmt.Println("异常结束:", err)
    } else {
        fmt.Println("正常结束,结果:", retry)
    }
}
