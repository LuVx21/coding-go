package api

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/lmittmann/tint"
	"gitlab.com/greyxor/slogor"
)

func Test_log_00(t *testing.T) {
	handler1 := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
			case slog.LevelKey:
				switch a.Value.String() {
				case "DEBUG":
				case "INFO":
				case "WARN":
				case "ERROR":
				}
			case slog.TimeKey:
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime + ".9999"))
				}
			}
			return a
		},
	})
	handler2 := slogor.NewHandler(os.Stderr,
		slogor.SetTimeFormat(time.DateTime+".9999"),
		slogor.SetLevel(slog.LevelDebug),
		slogor.ShowSource(),
	)
	handler3 := tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.DateTime + ".9999",
		Level:      slog.LevelDebug,
		AddSource:  true,
	})

	for _, logger := range []*slog.Logger{slog.Default(), slog.New(handler1), slog.New(handler2), slog.New(handler3)} {
		// logger := slog.New(h)
		// slog.SetDefault(logger)

		logger.Debug("这是一条调试信息")
		logger.Info("这是一条普通信息")
		logger.Warn("这是一条警告信息")
		logger.Error("这是一条错误信息")
	}
}
