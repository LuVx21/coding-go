package main

import (
    "fmt"
    "github.com/google/uuid"
    "time"

    "github.com/go-co-op/gocron/v2"
)

func main() {
    s, _ := gocron.NewScheduler()
    defer func() { _ = s.Shutdown() }()

    _, _ = s.NewJob(
        //gocron.DurationJob(
        //    10*time.Second,
        //),
        gocron.CronJob(
            "1/5 * * * * *",
            true,
        ),
        gocron.NewTask(
            func(a string, b int) {
                fmt.Println("执行任务开始->", time.Now(), "参数:", a, b)
            },
            "hello",
            1,
        ),
        gocron.WithName("测试任务"),
        gocron.WithEventListeners(
            gocron.BeforeJobRuns(
                func(jobID uuid.UUID, jobName string) {
                },
            ),
            gocron.AfterJobRuns(
                func(jobID uuid.UUID, jobName string) {
                    fmt.Println("执行任务完成-> id:", jobID, "任务名称:", jobName)
                },
            ),
            gocron.AfterJobRunsWithError(
                func(jobID uuid.UUID, jobName string, err error) {
                },
            ),
        ),
    )
    //fmt.Println(j.ID())

    s.Start()

    // block until you are ready to shut down
    select {
    case <-time.After(time.Minute):
    }

    _ = s.Shutdown()
}
