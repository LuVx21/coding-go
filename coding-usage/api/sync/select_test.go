package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_select(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 向通道1发送数据的goroutine
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Message from channel 1"
	}()

	// 向通道2发送数据的goroutine
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "Message from channel 2"
	}()

	// 使用select语句监听多个通道
	select {
	case msg1 := <-ch1:
		fmt.Println("Received from channel 1:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received from channel 2:", msg2)
	case <-time.After(3 * time.Second): // 超时控制
		fmt.Println("Timeout! No message received within 3 seconds.")
	}
}
