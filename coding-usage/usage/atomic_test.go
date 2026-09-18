package main

import (
	"fmt"
	"testing"

	"go.uber.org/atomic"
)

func Test_atomic_00(t *testing.T) {
	var message atomic.String

	// 设置值
	message.Store("hello")

	// 读取值
	msg := message.Load()
	fmt.Println("Message:", msg)

	// 交换值
	old := message.Swap("world")
	fmt.Println("Old:", old, "->", "New:", message.Load())
}
