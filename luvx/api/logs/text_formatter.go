package logs

import (
    "bytes"
    "fmt"
    "github.com/sirupsen/logrus"
    "path"
    "time"
)

const (
    red    = 31
    yellow = 33
    blue   = 36
    gray   = 37
)

type MyLogFormatter struct {
    TimestampFormat string
}

func (t MyLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    //根据不同的level展示颜色
    var levelColor int
    switch entry.Level {
    case logrus.DebugLevel, logrus.TraceLevel:
        levelColor = gray
    case logrus.WarnLevel:
        levelColor = yellow
    case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
        levelColor = red
    default:
        levelColor = blue
    }
    //字节缓冲区
    var b *bytes.Buffer
    if entry.Buffer != nil {
        b = entry.Buffer
    } else {
        b = &bytes.Buffer{}
    }
    //自定义日期格式
    timestampFormat := t.TimestampFormat
    if timestampFormat == "" {
        timestampFormat = time.RFC3339
    }
    timestamp := entry.Time.Format(timestampFormat)
    if entry.HasCaller() {
        fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
        funcVal := entry.Caller.Function
        _, _ = fmt.Fprintf(b, "%s \033[%dm%s\033[0m %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
    } else {
        _, _ = fmt.Fprintf(b, "%s \033[%dm%s\033[0m %s\n", timestamp, levelColor, entry.Level, entry.Message)
    }
    return b.Bytes(), nil
}
