package main

import (
    "fmt"
    "github.com/loov/hrtime"
    "time"
)

func main() {
    start := hrtime.Now()
    time.Sleep(time.Second)
    fmt.Println("耗时:", hrtime.Since(start))

    const numberOfExperiments = 4096
    bench := hrtime.NewBenchmark(numberOfExperiments)
    for bench.Next() {
        time.Sleep(1000 * time.Nanosecond)
    }
    fmt.Println(bench.Histogram(10))
}
