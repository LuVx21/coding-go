package goroutine

import (
	"fmt"
	"sync"
	"testing"
)

func Test_00(t *testing.T) {
	fmt.Println("main", GoID())
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Go(func() {
			fmt.Println(i, GoID())
		})
	}
	wg.Wait()
}
