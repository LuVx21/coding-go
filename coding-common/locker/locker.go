package locker

import (
	"time"
)

// Locker 使用例: var l locker.Locker[string] = &EtcdLocker[string]{}
type Locker[T any] interface {
	// TryLock 加锁
	TryLock(key T, exp time.Duration) bool
	// Unlock 解锁
	Unlock(key T) bool
}
