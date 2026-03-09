package syncx

import (
	"sync"
	"time"
)

// TryLockWithTimeout 超时后获取锁失败
func TryLockWithTimeout(mu *sync.Mutex, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		mu.Lock()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}
