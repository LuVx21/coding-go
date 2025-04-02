package common_x

import (
	"log/slog"
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
func IfThenGet[T any](expr bool, a, b func() T) T {
	if expr {
		return a()
	}
	return b()
}

// RunCatching 捕捉异常,避免异常退出
func RunCatching(fn func()) {
	RunCatchingReturn(func() int {
		fn()
		return 0
	})
}

func RunCatchingReturn[T any](fn func() T) T {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("fast-fail", r)
		}
	}()
	return fn()
}

func RunInRoutine(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}

func RunWithTime[R any](name string, f func() R) R {
	defer TrackTime1(name, time.Now())
	return f()
}

func RunWithTime2[R1 any, R2 any](name string, f func() (R1, R2)) (R1, R2) {
	defer TrackTime1(name, time.Now())
	return f()
}

func TrackTime1(name string, start time.Time) {
	elapsed := time.Since(start)
	slog.Info(name, "执行时间", elapsed)
}

func TrackTime(name string) func() {
	start := time.Now()
	return func() {
		slog.Info(name, "执行时间", time.Since(start))
	}
}
