package slogs

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/luvx21/coding-go/coding-common/os_x"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	infoLogFile  = "app.log"
	errorLogFile = "error.log"
)

var (
	logDir        = os_x.Getenv("HOME") + "/data/slogs"
	defaultLevel  = slog.LevelDebug
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
	infoHandler := slog.NewJSONHandler(infoLog, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	errorHandler := slog.NewJSONHandler(errorLog, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	defaultLogger = slog.New(newMultiHandler(consoleHandler, infoHandler, errorHandler))
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
