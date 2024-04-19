package runner

import (
    "github.com/go-co-op/gocron/v2"
    "github.com/google/uuid"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/service"
    "time"
)

func Start() {
    go exec()
}

func exec() {
    s, _ := gocron.NewScheduler()
    //defer func() { _ = s.Shutdown() }()

    _, _ = s.NewJob(
        gocron.CronJob(
            "0 7/10 * * * *",
            true,
        ),
        gocron.NewTask(
            service.PullHotBand,
        ),
        gocron.WithName("拉取微博热搜"),
        gocron.WithEventListeners(
            gocron.BeforeJobRuns(
                func(jobID uuid.UUID, jobName string) {
                    logs.Log.Infoln("执行任务开始-> id:", jobID, "任务名称:", jobName)
                },
            ),
            gocron.AfterJobRuns(
                func(jobID uuid.UUID, jobName string) {
                    logs.Log.Infoln("执行任务完成-> id:", jobID, "任务名称:", jobName)
                },
            ),
        ),
    )

    s.Start()

    select {
    case <-time.After(time.Minute):
    }
}
