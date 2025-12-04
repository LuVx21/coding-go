package slogs

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/luvx21/coding-go/coding-common/configs_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/infra/logs"
	"github.com/spf13/viper"
	"gitlab.com/greyxor/slogor"
)

var (
	defaultLevel              = slog.LevelInfo
	logDir                    = os_x.Getenv("HOME") + "/data/slogs"
	infoLogFile, errorLogFile = "main-slog", "error-slog"
	logFormat                 = "json"

	hs = []slog.Handler{}

	defaultLogger *slog.Logger
	initOnce      sync.Once
)

// RegisterHandler 添加自定义Handler
func RegisterHandler(h ...slog.Handler) {
	hs = append(hs, h...)
}

func SetConsoleLevel(level slog.Level) {
	defaultLevel = level
}

func SetLogDir(path string) {
	logDir = path
}

func initLogger() {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	infoWriter, errWriter := logs.LogWriter(logDir, infoLogFile), logs.LogWriter(logDir, errorLogFile)

	consoleHandler := slogor.NewHandler(os.Stderr, slogor.SetTimeFormat(time.DateTime+".9999"), slogor.SetLevel(defaultLevel), slogor.ShowSource())
	RegisterHandler(consoleHandler)
	if logFormat == "text" {
		infoHandler := slogor.NewHandler(infoWriter, slogor.SetTimeFormat(time.DateTime+".9999"), slogor.SetLevel(slog.LevelInfo), slogor.ShowSource(), slogor.DisableColor())
		errorHandler := slogor.NewHandler(errWriter, slogor.SetTimeFormat(time.DateTime+".9999"), slogor.SetLevel(slog.LevelError), slogor.ShowSource(), slogor.DisableColor())

		RegisterHandler(infoHandler, errorHandler)
	} else {
		infoHandler := slog.NewJSONHandler(infoWriter, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
		errorHandler := slog.NewJSONHandler(errWriter, &slog.HandlerOptions{Level: slog.LevelError, AddSource: true})

		RegisterHandler(infoHandler, errorHandler)
	}
	defaultLogger = slog.New(newMultiHandler(hs...))

	slog.SetDefault(defaultLogger)
}

// GetLogger returns the default logger
func GetLogger() *slog.Logger {
	initOnce.Do(initLogger)
	return defaultLogger
}

func InitFromConfig(c *viper.Viper) {
	if c == nil {
		c = configs_x.GetConfigByKey("log")
	}
	var lc logs.LogConfig
	if c != nil && c.Unmarshal(&lc) == nil {
		if len(lc.Level) > 0 {
			l := slog.LevelInfo
			switch strings.ToUpper(lc.Level) {
			case "DEBUG":
				l = slog.LevelDebug
			case "WARN":
				l = slog.LevelWarn
			case "ERROR":
				l = slog.LevelError
			}
			if l != slog.LevelInfo {
				SetConsoleLevel(l)
			}
		}
		if len(lc.LogDir) > 0 {
			if strings.HasPrefix(lc.LogDir, "/") || strings.HasPrefix(lc.LogDir, "$") {
				SetLogDir(os.ExpandEnv(lc.LogDir))
			} else {
				if exePath, err := os.Executable(); err == nil {
					SetLogDir(filepath.Join(filepath.Dir(exePath), lc.LogDir))
				}
			}
		}
		if len(lc.MainLog) > 0 {
			infoLogFile = lc.MainLog
		}
		if len(lc.ErrorLog) > 0 {
			errorLogFile = lc.ErrorLog
		}
		if len(lc.LogFormat) > 0 && (lc.LogFormat == "text" || lc.LogFormat == "json") {
			logFormat = lc.LogFormat
		}
	}

	GetLogger()
}
