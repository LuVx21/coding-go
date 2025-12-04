package slogs

import (
	"context"
	"log/slog"
)

// 可使用 https://github.com/samber/slog-multi

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
