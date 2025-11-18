package runner

import (
	"context"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/db"
	"luvx/gin/runner/docker"
	"luvx/gin/service"
	"luvx/gin/service/bili"
	"luvx/gin/service/keeplive"
	"luvx/gin/service/rss"
	"luvx/gin/service/weibo_p"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/luvx21/coding-go/coding-common/cast_x"
)

var (
	beforeListener = gocron.BeforeJobRuns(
		func(jobID uuid.UUID, jobName string) {
			// logs.Log.Infoln("任务:", jobName, "开始")
		},
	)
	afterListener = gocron.AfterJobRuns(
		func(jobID uuid.UUID, jobName string) {
			// logs.Log.Infoln("任务:", jobName, "完成")
		},
	)
)

func Start() {
	result, _ := db.RedisClient.HGet(context.TODO(), consts.AppSwitchKey, "runner_all").Result()
	if !cast_x.ToBool(result) {
		return
	}
	go exec()
}

func exec() {
	s, _ := gocron.NewScheduler()
	// defer func() { _ = s.Shutdown() }()

	callRunnerRegister(s)

	s.Start()

	time.Sleep(time.Minute)

	// select {
	// case <-time.After(time.Minute):
	// }
}

func callRunnerRegister(s gocron.Scheduler) {
	var runners []*service.Runner
	runners = append(runners, weibo_p.RunnerRegister()...)
	runners = append(runners, rss.RunnerRegister()...)
	runners = append(runners, bili.RunnerRegister()...)
	runners = append(runners, keeplive.RunnerRegister()...)
	runners = append(runners, docker.RunnerRegister()...)
	for _, r := range runners {
		_, _ = s.NewJob(
			gocron.CronJob(r.Crontab, true),
			gocron.NewTask(r.Fn),
			gocron.WithName(r.Name),
			gocron.WithEventListeners(beforeListener, afterListener),
		)
	}
}
