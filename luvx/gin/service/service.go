package service

import (
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/infra/infra_sql"
	"luvx/gin/db"
)

var db1, _ = db.MySQLClient.DB()
var RunnerLocker = infra_sql.NewLocker[string](db1)

type Runner struct {
	Name    string
	Crontab string
	Fn      func()
}

func NewRunner(key, crontab string, exp time.Duration, f func()) *Runner {
	fn := func() { RunnerLocker.LockRun(key, exp, f) }
	return &Runner{Name: key, Crontab: crontab, Fn: func() { common_x.RunCatching(fn) }}
}
