package slogs

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/infra/logs"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultLevel              = slog.LevelInfo
	logDir                    = os_x.Getenv("HOME") + "/data/slogs"
	infoLogFile, errorLogFile = "app.log", "error.log"
	logFormat                 = "json"

	defaultLogger *slog.Logger
	initOnce      sync.Once
)

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

	infoLog := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, infoLogFile),
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	errorLog := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, errorLogFile),
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     defaultLevel,
	})

	if logFormat == "text" {
		infoHandler := slog.NewTextHandler(infoLog, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		errorHandler := slog.NewTextHandler(errorLog, &slog.HandlerOptions{
			Level: slog.LevelError,
		})
		defaultLogger = slog.New(newMultiHandler(consoleHandler, infoHandler, errorHandler))
	} else {
		infoHandler := slog.NewJSONHandler(infoLog, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		errorHandler := slog.NewJSONHandler(errorLog, &slog.HandlerOptions{
			Level: slog.LevelError,
		})
		defaultLogger = slog.New(newMultiHandler(consoleHandler, infoHandler, errorHandler))
	}

	slog.SetDefault(defaultLogger)
}

type multiHandler struct {
	handlers []slog.Handler
}

func newMultiHandler(handlers ...slog.Handler) *multiHandler {
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var handlers []slog.Handler
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithAttrs(attrs))
	}
	return newMultiHandler(handlers...)
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	var handlers []slog.Handler
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithGroup(name))
	}
	return newMultiHandler(handlers...)
}

// GetLogger returns the default logger
func GetLogger() *slog.Logger {
	initOnce.Do(initLogger)
	return defaultLogger
}

type temp struct {
	Log logs.LogConfig
}

func InitFromConfig(c *viper.Viper) {
	var t temp
	c.Unmarshal(&t)

	if len(t.Log.Level) > 0 {
		l := slog.LevelInfo
		switch strings.ToUpper(t.Log.Level) {
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
	if len(t.Log.LogDir) > 0 {
		if strings.HasPrefix(t.Log.LogDir, "/") || strings.HasPrefix(t.Log.LogDir, "$") {
			SetLogDir(os.ExpandEnv(t.Log.LogDir))
		} else {
			if exePath, err := os.Executable(); err == nil {
				SetLogDir(filepath.Join(filepath.Dir(exePath), t.Log.LogDir))
			}
		}
	}
	if len(t.Log.MainLog) > 0 {
		infoLogFile = t.Log.MainLog
	}
	if len(t.Log.ErrorLog) > 0 {
		errorLogFile = t.Log.ErrorLog
	}
	if len(t.Log.LogFormat) > 0 && (t.Log.LogFormat == "text" || t.Log.LogFormat == "json") {
		logFormat = t.Log.LogFormat
	}

	GetLogger()
}
