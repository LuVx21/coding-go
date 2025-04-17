package syncx

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// 带超时控制的请求合并
func DoWithTimeout(group *singleflight.Group, key string, fn func() (any, error), timeout time.Duration) (any, error) {
	result, errChan := make(chan any, 1), make(chan error, 1)
	go func() {
		res, err, _ := group.Do(key, fn)
		if err != nil {
			errChan <- err
			return
		}
		result <- res
	}()

	select {
	case res := <-result:
		return res, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for result")
	}
}

type CachedGroup struct {
	group singleflight.Group
	cache sync.Map
	ttl   time.Duration
}

// 结果缓存扩展
func (g *CachedGroup) Do(key string, fn func() (any, error)) (any, error) {
	if val, ok := g.cache.Load(key); ok {
		return val, nil
	}

	val, err, _ := g.group.Do(key, func() (any, error) {
		res, err := fn()
		if err == nil {
			time.AfterFunc(g.ttl, func() { g.cache.Delete(key) })
		}
		return res, err
	})

	if err == nil {
		g.cache.Store(key, val)
	}
	return val, err
}
