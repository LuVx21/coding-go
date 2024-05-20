package goroutine

import (
    "fmt"
    "sync"
    "testing"
)

func Test_00(t *testing.T) {
    fmt.Println("main", GoID())
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Println(i, GoID())
        }()
    }
    wg.Wait()
}
