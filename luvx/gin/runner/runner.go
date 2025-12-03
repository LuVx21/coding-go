package runner

import (
	"context"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/db"
	"luvx/gin/service"
	"luvx/gin/service/bili"
	"luvx/gin/service/keeplive"
	"luvx/gin/service/rss"
	"luvx/gin/service/weibo_p"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	RunnerMap = make(map[string]func(), 16)
)

func Start() {
	result, err := db.RedisClient.HGet(context.TODO(), consts.AppSwitchKey, "runner_all").Bool()
	if err != nil || !result {
		logrus.Warn("定时任务未启用", err, result)
		return
	}
	go exec()
}

func exec() {
	s, _ := gocron.NewScheduler()
	// defer func() { _ = s.Shutdown() }()

	callRunnerRegister(s)

	s.Start()

	// time.Sleep(time.Minute)

	select {
	case <-time.After(time.Minute):
	}
}

func callRunnerRegister(s gocron.Scheduler) {
	var runners []*service.Runner
	runners = append(runners, weibo_p.RunnerRegister()...)
	runners = append(runners, rss.RunnerRegister()...)
	runners = append(runners, bili.RunnerRegister()...)
	runners = append(runners, keeplive.RunnerRegister()...)
	// runners = append(runners, xxx.RunnerRegister()...)
	for _, r := range runners {
		logrus.Infof("定时任务已配置 %-20s %s", r.Crontab, r.Name)
		RunnerMap[r.Name] = r.Fn
		_, _ = s.NewJob(
			gocron.CronJob(r.Crontab, true),
			gocron.NewTask(r.Fn),
			gocron.WithName(r.Name),
			gocron.WithEventListeners(beforeListener, afterListener),
		)
	}
}
