package runner

import (
    "github.com/go-co-op/gocron/v2"
    "github.com/google/uuid"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/service"
    "luvx/gin/service/rss"
    "luvx/gin/service/weibo_p"
    "time"
)

var (
    beforeListener = gocron.BeforeJobRuns(
        func(jobID uuid.UUID, jobName string) {
            logs.Log.Infoln("任务:", jobName, "开始")
        },
    )
    afterListener = gocron.AfterJobRuns(
        func(jobID uuid.UUID, jobName string) {
            logs.Log.Infoln("任务:", jobName, "完成")
        },
    )
)

func Start() {
    go exec()
}

func exec() {
    s, _ := gocron.NewScheduler()
    //defer func() { _ = s.Shutdown() }()

    callRunnerRegister(s)

    s.Start()

    select {
    case <-time.After(time.Minute):
    }
}

func callRunnerRegister(s gocron.Scheduler) {
    var runners []*service.Runner
    runners = append(runners, weibo_p.RunnerRegister()...)
    runners = append(runners, rss.RunnerRegister()...)
    for _, r := range runners {
        _, _ = s.NewJob(
            gocron.CronJob(r.Crontab, true),
            gocron.NewTask(r.Fn),
            gocron.WithName(r.Name),
            gocron.WithEventListeners(beforeListener, afterListener),
        )
    }
}
