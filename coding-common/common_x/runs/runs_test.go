package runs

import (
	"fmt"
	"testing"
	"time"
)

func Test_00(t *testing.T) {
	Go(a)

	time.Sleep(time.Second * 3)
}

func a() {
	fmt.Println("haha")
}
