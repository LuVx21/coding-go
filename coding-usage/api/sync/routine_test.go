package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

func Test_return(t *testing.T) {
	N := 2

	wg := sync.WaitGroup{}
	// 定义⼀个⽤于存放结果的channel
	resultCh := make(chan int, N)

	for i := 1; i <= N; i++ {
		wg.Add(1)
		// 启动协程计算结果，并将结果发送到channel中
		go func(num int) {
			defer wg.Done()
			time.Sleep(time.Second)
			resultCh <- num * 2
		}(i)
	}

	// 等待所有协程执⾏完毕
	go func() {
		fmt.Println("point1")
		wg.Wait()
		// 关闭channel
		close(resultCh)
		fmt.Println("point2")
	}()
	fmt.Println("point3")

	// 从channel中读取所有结果，并将结果累加起来
	sum := 0
	for res := range resultCh {
		fmt.Println("结果", res)
		sum += res
	}
	fmt.Println("和:", sum)
}

func Test_02(t *testing.T) {
	f1 := func() string {
		log.Print("task start")
		time.Sleep(time.Second * 1)
		log.Print("task done")
		return "ok"
	}
	f2 := func() int {
		log.Print("task start")
		time.Sleep(time.Second * 2)
		log.Print("task done")
		return 100
	}

	start := time.Now()
	fmt.Println(start)

	wg := sync.WaitGroup{}

	// 多种类型的结果,可以使用
	// make(chan any, numWorkers)
	r1 := make(chan string, 1)
	r2 := make(chan int, 1)

	wg.Go(func() { r1 <- f1() })
	wg.Go(func() { r2 <- f2() })

	//make方法,不指定缓冲区大小,这里直接使用会出现死锁
	wg.Wait()
	close(r1)
	close(r2)

	//make方法,不指定缓冲区大小时,可使用
	//go func() {
	//    wg.Wait()
	//    close(r1)
	//    close(r2)
	//}()

	s1 := <-r1
	s2 := <-r2
	fmt.Println(s1, s2)
	fmt.Println(time.Since(start))
}

func Test_03(t *testing.T) {
	tasks := []string{"task1", "task2", "task3", "task4", "task5", "task6", "task7"}

	// 创建并发控制通道，容量为3
	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup

	for i, task := range tasks {
		sem <- struct{}{} // 获取信号量
		wg.Go(func() {
			defer func() { <-sem }() // 释放信号量
			fmt.Printf("开始执行任务 %d: %s\n", i, task)
			time.Sleep(time.Second * time.Duration(common_x.IfThen(i%2 == 0, 2, 3))) // 模拟任务执行
			fmt.Printf("完成任务 %d: %s\n", i, task)
		})
	}

	wg.Wait()
	fmt.Println("所有任务执行完成")
}
