package main

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_once_00(t *testing.T) {
	var r string
	var once sync.Once
	for i := range 10 {
		go func(ii int) {
			s := strconv.Itoa(ii)
			fmt.Println("aaaaa" + " No." + s)
			once.Do(func() {
				r = func(iii int) string {
					fmt.Println("初始化: " + strconv.Itoa(iii))
					time.Sleep(time.Second * 2)
					return "初始化: " + time.Now().String() + " No." + strconv.Itoa(iii)
				}(ii)
			})
			fmt.Println("bbbbb" + " No." + s)
		}(i)
	}
	time.Sleep(time.Second * 5)
	fmt.Println(r)
}

func Test_once_01(t *testing.T) {
	rf := sync.OnceValue(func() string {
		fmt.Println("aaaaa")
		return "初始化: " + time.Now().String()
	})

	var r string
	for i := range 10 {
		fmt.Println(i)
		r = rf()
	}

	fmt.Println(r)

	time.Sleep(time.Second * 3)
}
