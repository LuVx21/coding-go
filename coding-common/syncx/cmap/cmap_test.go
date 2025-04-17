package cmap

import (
	"strconv"
	"sync"
	"testing"
)

type ThreadSafeMap interface {
	Get(key string) any
	Set(key string, val any)
}

type ConcurMap struct {
	m ConcurrentMap[string, any]
}

func (c *ConcurMap) Get(key string) any {
	v, _ := c.m.Get(key)
	return v
}
func (c *ConcurMap) Set(key string, val any) {
	c.m.Set(key, val)
}

func benchmark(b *testing.B, m ThreadSafeMap, read, write int) {
	for i := 0; b.Loop(); i++ {
		var wg sync.WaitGroup

		// 注意: 这里的读写操作有一部分 key 是重合的

		// 读操作
		for k := range read * 100 {
			wg.Add(1)
			go func(key int) {
				m.Get(strconv.Itoa(i * key))
				wg.Done()
			}(k)
		}

		// 写操作
		for k := range write * 100 {
			wg.Add(1)
			go func(key int) {
				m.Set(strconv.Itoa(i*key), key)
				wg.Done()
			}(k)
		}

		wg.Wait()
	}
}

// go test -count=1 -run='^$' -bench=. -benchtime=3s  -benchmem
func BenchmarkReadMoreRWMutex(b *testing.B) {
	benchmark(b, &ConcurMap{m: New[any]()}, 9, 1)
}
