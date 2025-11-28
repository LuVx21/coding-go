package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var r string

func onlyOneReturn() string {
	time.Sleep(time.Second * 2)
	fmt.Println("初始化")
	return "初始化: " + time.Now().String()
}
func Test_once_00(t *testing.T) {
	var once sync.Once
	for range 10 {
		go func() {
			fmt.Println("aaaaa")
			once.Do(func() {
				r = onlyOneReturn()
			})
			fmt.Println("bbbbb")
		}()
	}
	time.Sleep(time.Second * 5)
	fmt.Println(r)
}
