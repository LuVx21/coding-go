package runner

import (
    "github.com/go-co-op/gocron/v2"
    "github.com/google/uuid"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/service/weibo_p"
    "time"
)

func Start() {
    go exec()
}

func exec() {
    s, _ := gocron.NewScheduler()
    //defer func() { _ = s.Shutdown() }()

    beforeListener := gocron.BeforeJobRuns(
        func(jobID uuid.UUID, jobName string) {
            logs.Log.Infoln("执行任务开始-> id:", jobID, "任务名称:", jobName)
        },
    )
    afterListener := gocron.AfterJobRuns(
        func(jobID uuid.UUID, jobName string) {
            logs.Log.Infoln("执行任务完成-> id:", jobID, "任务名称:", jobName)
        },
    )
    _, _ = s.NewJob(
        gocron.CronJob("0 7/10 * * * *", true),
        gocron.NewTask(func() { common_x.RunCatching(weibo_p.PullHotBand) }),
        gocron.WithName("拉取微博热搜"),
        gocron.WithEventListeners(beforeListener, afterListener),
    )
    _, _ = s.NewJob(
        gocron.CronJob("0 9/10 * * * *", true),
        gocron.NewTask(func() { common_x.RunCatching(weibo_p.PullByGroup) }),
        gocron.WithName("拉取分组微博"),
        gocron.WithEventListeners(beforeListener, afterListener),
    )

    s.Start()

    select {
    case <-time.After(time.Minute):
    }
}
