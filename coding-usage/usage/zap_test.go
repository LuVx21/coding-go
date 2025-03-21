package main

import (
	"flag"
	"testing"
	"time"

	"github.com/golang/glog"
	"go.uber.org/zap"
)

func Test_b(t *testing.T) {
	glog.MaxSize = 1024 * 1024 * 1024 // 1G自动分割
	flag.Parse()
	defer glog.Flush()

	glog.Info("This is info message")
	glog.Infof("This is info message: %v", 123)

	glog.Warning("This is warning message")
	glog.Warningf("This is warning message: %v", 123)

	glog.Error("This is error message")
	glog.Errorf("This is error message: %v", 123)

	// glog.Fatal("This is fatal message")
	// glog.Fatalf("This is fatal message: %v", 123)
}

func Test_c(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	url := "http://marmotedu.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
