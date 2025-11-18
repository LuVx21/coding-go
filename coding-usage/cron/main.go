package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	m0()
}

func m0() {
	c := cron.New(cron.WithSeconds())
	defer c.Stop()

	// 每5秒执行一次
	_, _ = c.AddFunc("@every 5s", func() { fmt.Printf("每5秒, %s\n", time.Now().Format("15:04:05")) })
	// 每10秒执行一次
	_, _ = c.AddFunc("@every 10s", func() { fmt.Printf("每10秒, %s\n", time.Now().Format("15:04:05")) })
	_, _ = c.AddFunc("0 0/1 * * * *", func() { fmt.Printf("每1分, %s\n", time.Now().Format("15:04:05")) })
	go c.Start()

	select {}
}

func m1() {
	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	beforeListener := gocron.BeforeJobRuns(
		func(jobID uuid.UUID, jobName string) {
		},
	)
	afterListener := gocron.AfterJobRuns(
		func(jobID uuid.UUID, jobName string) {
			fmt.Println("执行任务完成-> id:", jobID, "任务名称:", jobName)
		},
	)
	withError := gocron.AfterJobRunsWithError(
		func(jobID uuid.UUID, jobName string, err error) {
		},
	)
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
		gocron.WithEventListeners(beforeListener, afterListener, withError),
	)
	//fmt.Println(j.ID())

	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute):
	}

	//_ = s.Shutdown()
}
