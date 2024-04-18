package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/sethvargo/go-retry"
)

func main() {
    f := func(ctx context.Context) error {
        if err := test(ctx, true); err != nil {
            return retry.RetryableError(err)
        }
        return nil
    }
    if err := retry.Fibonacci(context.Background(), 1*time.Second, f); err != nil {
        log.Fatal(err)
    }
}

func test(ctx context.Context, b bool) error {
    fmt.Println("执行方法......")
    err := fmt.Errorf("异常")
    if b {
        return err
    }
    return nil
}
