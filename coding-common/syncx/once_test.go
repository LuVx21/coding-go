package syncx

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test_lazy_one_00(t *testing.T) {
	ol := LazyOnceIn(func(ii int) string {
		fmt.Println("执行一次!!! No." + strconv.Itoa(ii))
		return "初始化: " + time.Now().String() + " No." + strconv.Itoa(ii)
	})

	results := make(chan string, 10)
	for i := range 10 {
		go func(ii int) {
			fmt.Println("aaaaa No." + strconv.Itoa(ii))
			r := ol(ii)
			fmt.Println("bbbbb No." + strconv.Itoa(ii))
			results <- r
		}(i)
	}
	time.Sleep(time.Second)
	ol(-1)

	time.Sleep(time.Second * 2)

	for len(results) > 0 {
		r := <-results
		fmt.Println("结果:", r)
	}
}
