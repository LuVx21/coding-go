package runs

import (
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func Defered(f func()) func() {
	f1 := func(_ int) { f() }
	return func() { DeferedArgs(f1)(-1) }
}

// DeferedArgs 函数外包装一层recover
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

func GracefulStop(before, after func()) {
	if before != nil {
		before()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	slog.Info("GracefulStop...", "sig", sig)

	if after != nil {
		after()
	}
}
