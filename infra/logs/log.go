package logs

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	prefixed "github.com/luvx12/logrus-prefixed-formatter"
	"github.com/luvx21/coding-go/coding-common/configs_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	defaultLevel              = logrus.InfoLevel
	logDir                    = os_x.Getenv("HOME") + "/data/logs"
	infoLogFile, errorLogFile = "main", "error"
	initOnce                  sync.Once
	// Log logger
	// Deprecated: 直接使用logrus
	Log = logrus.New()

	stdFormatter, fileFormatter *prefixed.TextFormatter // 命令行,文件输出格式
)

func SetConsoleLevel(level logrus.Level) {
	defaultLevel = level
}

func SetLogDir(path string) {
	logDir = path
}

func init() { initOnce.Do(initLogger) }
func initLogger() {
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
	dir, _ := os.Executable()
	logPath := path.Join(strings.TrimSpace(filepath.Dir(dir)), ".logs")
	logDir := os.Getenv("log_LogDir")
	if len(logDir) != 0 {
		logPath = logDir
	}

	writer, writer1 := LogWriter(logPath, infoLogFile), LogWriter(logPath, errorLogFile)

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

// GetLogger returns the default logger
func GetLogger() *logrus.Logger {
	initOnce.Do(initLogger)
	return Log
}

func InitFromConfig(c *viper.Viper) {
	if c == nil {
		c = configs_x.GetConfigByKey("log")
	}
	var lc LogConfig
	if c != nil && c.Unmarshal(&lc) == nil {
	}
	GetLogger()
}
