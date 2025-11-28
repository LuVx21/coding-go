package syncx

import "sync"

// LazyOnce 只执行一次的懒加载实现
func LazyOnce[T any](f func() T) func() T {
	var r T
	var once sync.Once
	return func() T {
		once.Do(func() { r = f() })
		return r
	}
}

// LazyOnce2 两个返回值
func LazyOnce2[A, B any](f func() (A, B)) func() (A, B) {
	var a A
	var b B
	var once sync.Once
	return func() (A, B) {
		once.Do(func() { a, b = f() })
		return a, b
	}
}
