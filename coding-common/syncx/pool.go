package syncx

import (
	"sync"
)

type pool[T any] struct {
	p sync.Pool
}

func NewPool[T any](fn func() T) *pool[T] {
	return &pool[T]{
		p: sync.Pool{New: func() any { return fn() }},
	}
}

func (p *pool[T]) Get() T  { return p.p.Get().(T) }
func (p *pool[T]) Put(t T) { p.p.Put(t) }
