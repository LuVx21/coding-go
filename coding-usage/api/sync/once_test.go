package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var r string

func Test_once_00(t *testing.T) {
	var once sync.Once
	for range 10 {
		go func() {
			fmt.Println("aaaaa")
			once.Do(func() {
				r = func() string {
					fmt.Println("初始化")
					time.Sleep(time.Second * 2)
					return "初始化: " + time.Now().String()
				}()
			})
			fmt.Println("bbbbb")
		}()
	}
	time.Sleep(time.Second * 5)
	fmt.Println(r)
}
