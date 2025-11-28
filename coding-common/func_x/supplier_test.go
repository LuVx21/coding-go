package func_x

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_er_00(t *testing.T) {
	er := Lazy(func() string { return time.Now().String() })
	fmt.Println(er.String())
	fmt.Println(er.Get())

	er.IfPresent(func(s string) error {
		fmt.Println("现在是", s)
		return nil
	})

	ss, _ := er.Map(func(s string) any {
		return "现在是:" + s
	})
	fmt.Println(ss)
	er.TryClose()
	fmt.Println(er.IsInitialized(), er.String())
}

func Test_er_01(t *testing.T) {
	callCount := 0
	supplier := Lazy(func() int {
		callCount++
		time.Sleep(10 * time.Millisecond)
		return 100
	})

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			value := supplier.Get()
			if value != 100 {
				t.Errorf("协程 %d: 期望 100, 得到 %d", i+1, value)
			}
		})
	}

	wg.Wait()

	// 验证懒加载函数只被调用一次
	if callCount != 1 {
		t.Errorf("期望懒加载函数被调用1次，实际调用 %d 次", callCount)
	}

	fmt.Printf("并发测试通过: 懒加载函数调用 %d 次\n", callCount)
}
