package infra_kv

import (
	"log/slog"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	infra_badger "github.com/luvx21/coding-go/infra/infra_kv/badgers"
)

type BadgerLocker[T any] struct {
	Client *badger.DB
	mu     sync.Mutex
}

func NewLocker[T any](c *badger.DB) *BadgerLocker[T] {
	return &BadgerLocker[T]{Client: c}
}

func (l *BadgerLocker[T]) TryLock(key T, exp time.Duration) bool {
	if !l.mu.TryLock() {
		return false
	}
	defer l.mu.Unlock()

	k := "lock:" + cast_x.ToString(key)
	_, exist := infra_badger.GetStr(l.Client, k)
	if exist {
		return false
	}

	err := infra_badger.SetStr(l.Client, k, "1", exp)
	return err == nil
}
func (l *BadgerLocker[T]) Unlock(key T) bool {
	if !l.mu.TryLock() {
		return false
	}
	defer l.mu.Unlock()

	k := "lock:" + cast_x.ToString(key)
	_, err := infra_badger.Delete(l.Client, []byte(k))
	return err == nil
}
func (l *BadgerLocker[T]) LockRun(key T, exp time.Duration, fn func()) {
	if l.TryLock(key, exp) {
		defer l.Unlock(key)
		fn()
	} else {
		slog.Warn("加锁失败", "key", key)
	}
}
