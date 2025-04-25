package etcds

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdLocker[T any] struct {
	Client *clientv3.Client
	mu     sync.Mutex
}

func NewLocker[T any](c *clientv3.Client) *EtcdLocker[T] {
	return &EtcdLocker[T]{Client: c}
}

func (l *EtcdLocker[T]) TryLock(key T, exp time.Duration) bool {
	if !l.mu.TryLock() {
		return false
	}
	defer l.mu.Unlock()

	k := "lock/" + cast_x.ToString(key)
	gr, err := l.Client.Get(context.Background(), k)
	if err != nil || (gr != nil && gr.Count > 0) {
		return false
	}
	err = Set(l.Client, k, "1", int64(exp.Seconds()))
	return err == nil
}
func (l *EtcdLocker[T]) Unlock(key T) bool {
	if !l.mu.TryLock() {
		return false
	}
	defer l.mu.Unlock()

	k := "lock/" + cast_x.ToString(key)
	_, err := l.Client.Delete(context.Background(), k)
	return err == nil
}
func (l *EtcdLocker[T]) LockRun(key T, exp time.Duration, fn func()) {
	if l.TryLock(key, exp) {
		defer l.Unlock(key)
		fn()
	} else {
		slog.Warn("加锁失败", "key", key)
	}
}
