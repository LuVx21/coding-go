package main

import (
    "fmt"
    "golang.org/x/sync/singleflight"
    "testing"
    "time"
)

var group singleflight.Group

func Test_00(t *testing.T) {
    go UsingSingleFlight("key")
    time.Sleep(1 * time.Second)
    go UsingSingleFlight("key")

    time.Sleep(2 * time.Second)

    go UsingSingleFlight("key")
    time.Sleep(1 * time.Second)
    go UsingSingleFlight("key")

    time.Sleep(2 * time.Second)
}

func UsingSingleFlight(key string) {
    v, err, shared := group.Do(key, func() (interface{}, error) {
        return FetchExpensiveData()
    })
    fmt.Println(v, err, shared)
}

//耗时或耗资源操作
func FetchExpensiveData() (int64, error) {
    fmt.Println("耗时或耗资源操作...", time.Now())
    time.Sleep(2 * time.Second)
    return time.Now().UnixNano(), nil
}
