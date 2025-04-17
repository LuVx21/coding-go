package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func initPrint() {
	time.Sleep(time.Second * 2)
	fmt.Println("初始化")
}
func Test_once_00(t *testing.T) {
	var once sync.Once
	for range 10 {
		go func() {
			fmt.Println("aaaaa")
			once.Do(initPrint)
			fmt.Println("bbbbb")
		}()
	}
	time.Sleep(time.Second * 5)
}
