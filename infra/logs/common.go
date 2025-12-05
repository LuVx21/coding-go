package logs

import (
	"io"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type LogConfig struct {
	Level             string
	LogDir            string
	MainLog, ErrorLog string // 日志文件名(不含扩展名)
	LogFormat         string
}

// LogWriter 日志文件轮转writer
func LogWriter(logPath, logName string) io.Writer {
	writer, _ := rotatelogs.New(
		path.Join(logPath, logName+"-%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join(logPath, logName+".log")),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	return writer
}
