package common

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "sync"
    "time"
)

func IfThen[T any](expr bool, a T, b T) T {
    if expr {
        return a
    }
    return b
}

// IfThenGet 使用时不简洁
func IfThenGet[T any](expr bool, a func() T, b func() T) T {
    if expr {
        return a()
    }
    return b()
}

// RunCatching 捕捉异常,避免异常退出
func RunCatching(fn func()) {
    func() {
        defer func() {
            if r := recover(); r != nil {
                log.Print(r)
            }
        }()
        fn()
    }()
}

func CatchingWithRoutine(fn func()) {
    go RunCatching(fn)
}

func RunInRoutine(wg *sync.WaitGroup, f func()) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        f()
    }()
}

func RunWithTime[R any](name string, f func() R) R {
    defer TrackTime(name, time.Now())
    return f()
}

func RunWithTime2[R1 any, R2 any](name string, f func() (R1, R2)) (R1, R2) {
    defer TrackTime(name, time.Now())
    return f()
}

func TrackTime(name string, start time.Time) {
    elapsed := time.Since(start)
    log.Infoln(name, "执行时间:", elapsed)
}

func TrackTime1(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s 执行时间: %v\n", name, time.Since(start))
    }
}