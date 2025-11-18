package main

import (
	"log"
	"time"
)

func Task() {
	// 模拟耗时任务
	log.Print("task start")
	time.Sleep(time.Second * 3)
	log.Print("task done")
}

func main() {
}
