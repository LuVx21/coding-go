package app

import (
	"log/slog"
	"runtime/debug"
	"time"

	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		{Name: "gc", Crontab: "0 7/15 * * * *", Fn: func() { common_x.RunCatching(gc) }},
	}
}

func gc() {
	service.RunnerLocker.LockRun("GC", time.Minute*3, func() {
		start := time.Now()
		// runtime.GC()
		debug.FreeOSMemory()
		elapsed := time.Since(start)
		slog.Debug("GC 完成", "耗时", elapsed)
	})
}
