package gorms

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

var GormLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	SlowThreshold:             300 * time.Millisecond,
	LogLevel:                  logger.Warn,
	IgnoreRecordNotFoundError: false,
	Colorful:                  true,
})
