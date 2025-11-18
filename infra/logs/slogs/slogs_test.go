package slogs

import (
	"log/slog"
	"testing"
)

func Test_01(t *testing.T) {
	SetConsoleLevel(slog.LevelInfo)
	logger := GetLogger()
	logger.Debug("这是一条调试信息")
	logger.Info("这是一条普通信息")
	logger.Warn("这是一条警告信息")
	logger.Error("这是一条错误信息")
}
