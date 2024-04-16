package main

import (
    "fmt"
    "strconv"
    "testing"
    "time"
)

func Test02(t *testing.T) {
    start := time.Now()
    ch := make(chan bool)
    for i := 0; i < 10; i++ {
        go parseUrls("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch)
    }

    for i := 0; i < 10; i++ {
        <-ch
    }

    elapsed := time.Since(start)
    fmt.Printf("Took %s", elapsed)
}
