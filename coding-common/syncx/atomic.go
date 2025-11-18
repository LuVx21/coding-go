package syncx

import "sync/atomic"

type AtomicValue[V any] struct {
	a atomic.Value
}

func NewAtomicValue[V any](v V) *AtomicValue[V] {
	a := new(AtomicValue[V])
	a.Store(v)
	return a
}

func (a *AtomicValue[V]) Store(val V) {
	a.a.Store(val)
}
func (a *AtomicValue[V]) Load() V {
	return a.a.Load().(V)
}
