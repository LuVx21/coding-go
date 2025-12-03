package logs

import (
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	prefixed "github.com/luvx12/logrus-prefixed-formatter"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	LogDir   string
	MainLog  string
	ErrorLog string
}

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

	// 以下三个常量都可以使用配置
	dir, _ := os_x.Command("sh", "-c", "go list -m -f {{.Dir}}")
	logPath := path.Join(strings.TrimSpace(dir), ".logs")
	logDir := os.Getenv("log_LogDir")
	if len(logDir) != 0 {
		logPath = logDir
	}
	writer, _ := rotatelogs.New(
		path.Join(logPath, "main-%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join(logPath, "main.log")),
		rotatelogs.WithMaxAge(time.Duration(168)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	writer1, _ := rotatelogs.New(
		path.Join(logPath, "error-%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join(logPath, "error.log")),
		rotatelogs.WithMaxAge(time.Duration(168)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
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
	Log.SetLevel(logrus.InfoLevel)

	logrus.AddHook(lfHook)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(stdFormatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	Log.Infoln("日志文件位置:", logPath)
}
