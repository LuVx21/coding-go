package main

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

func Test_00(t *testing.T) {
	go UsingSingleFlight("key")
	time.Sleep(1 * time.Second)
	go UsingSingleFlight("key")

	time.Sleep(2 * time.Second)

	go UsingSingleFlight("key")
	time.Sleep(1 * time.Second)
	go UsingSingleFlight("key")

	time.Sleep(2 * time.Second)
}

func UsingSingleFlight(key string) {
	v, err, shared := group.Do(key, FetchExpensiveData)
	fmt.Println(v, err, shared)
}

// 耗时或耗资源操作
func FetchExpensiveData() (any, error) {
	fmt.Println("耗时或耗资源操作...", time.Now())
	time.Sleep(2 * time.Second)
	return time.Now().UnixNano(), nil
}
func Test_singleflight_01(t *testing.T) {
	for range 10 {
		go UsingSingleFlight("key")
	}

	time.Sleep(time.Second * 10)
}
