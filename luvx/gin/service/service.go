package service

import (
	"context"
	"time"

	"luvx/gin/db"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/func_x"
	"github.com/luvx21/coding-go/infra/infra_sql"
	"go.mongodb.org/mongo-driver/bson"
)

var db1, _ = db.MySQLClient.DB()
var (
	RunnerLocker  = infra_sql.NewLocker[string](db1)
	DynamicConfig = func_x.Lazy(func() bson.M {
		var m bson.M
		db.GetCollection("config").FindOne(context.TODO(), bson.M{"_id": "app_config"}).Decode(&m)
		return m
	})
	DynamicCache = func_x.Lazy(func() bson.M {
		var m bson.M
		db.GetCollection("config").FindOne(context.TODO(), bson.M{"_id": "app_cache"}).Decode(&m)
		return m
	})
)

type Runner struct {
	Name, Crontab string
	Fn            func()
}

func NewRunner(key, crontab string, exp time.Duration, f func()) *Runner {
	fn := func() { RunnerLocker.LockRun(key, exp, f) }
	return &Runner{Name: key, Crontab: crontab, Fn: func() { common_x.RunCatching(fn) }}
}
