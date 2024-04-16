package main

import (
    "log"
    "sync"
    "testing"
)

func Test_01(t *testing.T) {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        Task()
        wg.Done()
    }()
    log.Print("111")
    wg.Wait()
    log.Print("main done")
}
