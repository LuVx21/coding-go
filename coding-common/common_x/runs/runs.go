package runs

import (
	"log/slog"
	"runtime/debug"
)

func Defered(f func()) func() {
	f1 := func(_ int) { f() }
	return func() { DeferedArgs(f1)(-1) }
}

func DeferedArgs[T any](f func(T)) func(T) {
	return func(t T) {
		defer func() {
			if err := recover(); err != nil {
				slog.Warn("异常", "panic", err)
				slog.Warn("异常栈", "stack", string(debug.Stack()))
			}
		}()

		f(t)
	}
}

// Go 野协程异常退出问题
func Go(f func()) {
	go Defered(f)()
}

// GoArgs 野协程异常退出问题(少见)
func GoArgs[T any](t T, f func(T)) {
	go DeferedArgs(f)(t)
}
