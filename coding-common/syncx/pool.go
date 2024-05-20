package syncx

import (
    "sync"
)

type Pool[T any] struct {
    pool sync.Pool
}

func NewPool[T any](fn func() T) Pool[T] {
    ff := func() any { return fn() }
    return Pool[T]{sync.Pool{New: ff}}
}

func (p *Pool[T]) Get() T {
    return p.pool.Get().(T)
}
func (p *Pool[T]) Put(x T) {
    p.pool.Put(x)
}
