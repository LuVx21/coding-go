package main

import (
    "context"
    "fmt"
    "github.com/luvx21/coding-go/coding-common/times_x"
    "golang.org/x/time/rate"
)

func main() {
    limiter := rate.NewLimiter(1, 1)
    for i := 0; i < 10; i++ {
        if err := limiter.Wait(context.Background()); err != nil {
            fmt.Println("Error waiting for limiter:", err)
            return
        }
        fmt.Println(times_x.TimeNow(), i+1)
    }
}
