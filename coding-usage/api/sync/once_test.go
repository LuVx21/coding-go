package main

import (
    "fmt"
    "sync"
    "testing"
    "time"
)

func initPrint() {
    fmt.Println("初始化")
}
func Test_once_00(t *testing.T) {
    var once sync.Once
    for i := 0; i < 10; i++ {
        go func(i int) {
            once.Do(initPrint)
        }(i)
    }
    time.Sleep(time.Second * 5)
}
