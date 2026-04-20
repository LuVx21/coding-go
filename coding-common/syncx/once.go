package syncx

import (
	"sync"
)

// LazyOnce 懒加载: `sync.Once`实现, 返回一个值
// `sync.OnceValue`的简化版, 无返回值的可用`sync.OnceFunc`
func LazyOnce[R any](f func() R) func() R {
	var r R
	var once sync.Once
	return func() R {
		once.Do(func() { r = f() })
		return r
	}
}

// LazyOnce 懒加载: `sync.Once`实现, 带入参, 返回一个值
func LazyOnceIn[I, R any](f func(I) R) func(I) R {
	var r R
	var once sync.Once
	return func(i I) R {
		once.Do(func() { r = f(i) })
		return r
	}
}
