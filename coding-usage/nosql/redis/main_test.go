package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func beforeAfter(caseName string) func() {
	if rdb == nil {
		rdb = connect()
	}
	return func() {
		fmt.Println(caseName, "end...")
	}
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()
	val, _ := rdb.Get(context.TODO(), "foo").Result()
	fmt.Println("foo", "=", val)

	r := rdb.Del(context.TODO(), "lock_"+"01")
	fmt.Println(r)

	b := rdb.SetNX(context.TODO(), "lock_"+"01", 1, time.Second*60)
	fmt.Println(b)
}

func Test_map_00(t *testing.T) {
	defer beforeAfter("Test_map_00")()
	result, _ := rdb.HGetAll(context.TODO(), "mm").Result()
	fmt.Println("mm", "=", result, result["mk"])

	v, _ := rdb.HGet(context.TODO(), "mm", "mk").Result()
	fmt.Println("mm.mk", "=", v)
}
