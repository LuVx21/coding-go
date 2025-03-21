package logs

import (
	"github.com/luvx21/coding-go/coding-common/consts_x"
	"github.com/sirupsen/logrus"
)

type TraceIdHook struct {
	TraceId string
}

func NewTraceIdHook(traceId string) logrus.Hook {
	return &TraceIdHook{
		TraceId: traceId,
	}
}

func (hook *TraceIdHook) Fire(entry *logrus.Entry) error {
	entry.Data[consts_x.TraceId] = hook.TraceId
	return nil
}

func (hook *TraceIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
