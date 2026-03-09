package syncx

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/times_x"
)

func Test_lock_00(t *testing.T) {
	var mu sync.Mutex
	// 不加锁, 下面就不超时,返回true
	// mu.Lock()
	defer mu.Unlock()

	fmt.Println(times_x.TimeNowDateSecond())
	b := TryLockWithTimeout(&mu, time.Second*3)
	fmt.Println(b)
	fmt.Println(times_x.TimeNowDateSecond())
}
