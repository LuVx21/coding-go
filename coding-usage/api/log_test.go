package api

import (
    "github.com/lmittmann/tint"
    "gitlab.com/greyxor/slogor"
    "log"
    "log/slog"
    "os"
    "testing"
    "time"
)

func Test_log_00(t *testing.T) {
    var handler slog.Handler
    handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelInfo,
        AddSource: true,
    })
    handler = slogor.NewHandler(os.Stderr, slogor.Options{
        TimeFormat: time.DateTime + ".99999",
        Level:      slog.LevelInfo,
        ShowSource: false,
    })
    handler = tint.NewHandler(os.Stderr, &tint.Options{
        TimeFormat: time.DateTime + ".99999",
        Level:      slog.LevelInfo,
        AddSource:  true,
    })

    logger := slog.New(handler)
    slog.SetDefault(logger)
    slog.Info("haha")

    log.Println("haha")
}

func Test_log_01(t *testing.T) {
}
