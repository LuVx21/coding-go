package app

import (
	"log/slog"
	"runtime/debug"
	"time"

	"luvx/gin/service"
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		service.NewRunner("gc", "0 7/15 * * * *", time.Minute*10, gc),
	}
}

func gc() {
	start := time.Now()
	// runtime.GC()
	debug.FreeOSMemory()
	elapsed := time.Since(start)
	slog.Debug("GC 完成", "耗时", elapsed)
}
