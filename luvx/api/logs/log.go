package logs

import (
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    prefixed "github.com/luvx12/logrus-prefixed-formatter"
    "github.com/rifflock/lfshook"
    "github.com/sirupsen/logrus"
    "luvx/config"
    "os"
    "path"
    "time"
)

var Log = logrus.New()

var stdFormatter *prefixed.TextFormatter  // 命令行输出格式
var fileFormatter *prefixed.TextFormatter // 文件输出格式

func init() {
    stdFormatter = &prefixed.TextFormatter{
        PrefixPadding:   3,
        FullTimestamp:   true,
        TimestampFormat: "2006-01-02 15:04:05.000000",
        ForceFormatting: true,
        ForceColors:     true,
        DisableColors:   false,
    }
    stdFormatter.SetColorScheme(&prefixed.ColorScheme{
        InfoLevelStyle:  "green",
        WarnLevelStyle:  "yellow",
        ErrorLevelStyle: "red",
        FatalLevelStyle: "41",
        PanicLevelStyle: "41",
        DebugLevelStyle: "blue",
        PrefixStyle:     "cyan",
        TimestampStyle:  "37",
        MessageStyle:    "37",
    })

    fileFormatter = &prefixed.TextFormatter{
        FullTimestamp:   true,
        TimestampFormat: "2006-01-02 15:04:05.000000",
        ForceFormatting: true,
        ForceColors:     false,
        DisableColors:   true,
    }

    logPath := config.AppConfig.Log.LogDir
    writer, _ := rotatelogs.New(
        path.Join(logPath, config.AppConfig.Log.MainLog),
        rotatelogs.WithLinkName(path.Join(logPath, "main.log")),
        rotatelogs.WithMaxAge(time.Duration(168)*time.Second),
        rotatelogs.WithRotationTime(time.Duration(24)*time.Second),
    )
    writer1, _ := rotatelogs.New(
        path.Join(logPath, config.AppConfig.Log.ErrorLog),
        rotatelogs.WithLinkName(path.Join(logPath, "error.log")),
        rotatelogs.WithMaxAge(time.Duration(168)*time.Second),
        rotatelogs.WithRotationTime(time.Duration(24)*time.Second),
    )

    lfHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer,
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer1,
        logrus.FatalLevel: writer1,
        logrus.PanicLevel: writer1,
    }, fileFormatter)

    Log.AddHook(lfHook)
    Log.SetReportCaller(true)
    Log.SetFormatter(stdFormatter)
    Log.SetOutput(os.Stdout)
    Log.SetLevel(logrus.DebugLevel)
}
