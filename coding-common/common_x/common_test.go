package common_x

import (
	"fmt"
	"log/slog"
	"testing"
	"time"
)

func Test_m1(t *testing.T) {
	RunCatching(func() {
		panic("异常")
	})
	fmt.Println("后续操作1")
	r := RunCatchingReturn(func() string {
		panic("异常")
		//return "结果"
	})
	fmt.Println(r)
	fmt.Println("后续操作2")
}

func Test_RunWithTime(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	withTime := RunWithTimeReturn("m1", func() string {
		time.Sleep(time.Second)
		return "ok"
	})

	fmt.Println(withTime)
}

func Test_01(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	defer TrackTime1("main", time.Now())
	defer TrackTime("main1")()
}
