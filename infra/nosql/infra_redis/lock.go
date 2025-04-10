package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/redis/go-redis/v9"
)

type Locker[T any] struct {
	Client *redis.Client
}

func NewLocker[T any](c *redis.Client) *Locker[T] {
	return &Locker[T]{Client: c}
}

func (l *Locker[T]) TryLock(key T, exp time.Duration) bool {
	r := l.Client.SetNX(context.TODO(), cast_x.ToString(key), 1, exp)
	return r.Val()
}
func (l *Locker[T]) Unlock(key T) bool {
	r := l.Client.Del(context.TODO(), cast_x.ToString(key))
	return r.Val() >= 0
}
func (l *Locker[T]) LockRun(key T, exp time.Duration, fn func()) {
	if l.TryLock(key, exp) {
		defer l.Unlock(key)
		fn()
	} else {
		slog.Warn("加锁失败", "key", key)
	}
}
