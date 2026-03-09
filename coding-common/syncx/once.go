package syncx

import "sync"

// Do 带返回值的泛型Once Do
func Do[T any](once *sync.Once, f func() T) T {
	var r T
	once.Do(func() {
		r = f()
	})
	return r
}

// LazyOnce 只执行一次的懒加载实现
func LazyOnce[T any](f func() T) func() T {
	var r T
	var once sync.Once
	return func() T {
		once.Do(func() { r = f() })
		return r
	}
}
