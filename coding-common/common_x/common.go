package common_x

import (
	"log/slog"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func IfThen[T any](expr bool, a T, b T) T {
	if expr {
		return a
	}
	return b
}

// IfThenGet 使用时不简洁
func IfThenGet[T any](expr bool, a, b func() T) T {
	return IfThen(expr, a, b)()
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
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			logrus.Errorln("fast-fail", "panic", r, "错误栈信息", string(buf[:n]))
		}
	}()
	return fn()
}

func RunWithTime(name string, f func()) {
	RunWithTimeReturn(name, func() int {
		f()
		return 0
	})
}

// RunWithTimeReturn 统计函数执行耗时
func RunWithTimeReturn[R any](name string, f func() R) R {
	defer TrackTime(name)()
	return f()
}

func TrackTime1(name string, start time.Time) {
	slog.Info("耗时统计", "统计项", name, "执行时间", time.Since(start))
}

func TrackTime(name string) func() {
	start := time.Now()
	return func() { TrackTime1(name, start) }
}
