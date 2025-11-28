package syncx

import (
	"fmt"
	"testing"
	"time"
)

func Test_lazy_one_00(t *testing.T) {
	ol := LazyOnce(func() string {
		fmt.Println("执行一次!!!")
		return time.Now().String()
	})

	for range 10 {
		go func() {
			fmt.Println("aaaaa")
			_ = ol()
			fmt.Println("bbbbb")
		}()
	}
	ol()
	time.Sleep(time.Second * 2)
}
