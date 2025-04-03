package locker

import (
	"sync"
	"time"
)

type Locker[T any] interface {
	// TryLock 加锁
	TryLock(key T, exp time.Duration) bool
	// Unlock 解锁
	Unlock(key T) bool
}

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
