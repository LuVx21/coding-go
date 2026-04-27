package syncx

import "sync/atomic"

type AtomicValue[V any] struct {
	v atomic.Value
}

func NewAtomicValue[V any](v V) *AtomicValue[V] {
	a := new(AtomicValue[V])
	a.Store(v)
	return a
}

func (a *AtomicValue[V]) Store(val V) { a.v.Store(val) }
func (a *AtomicValue[V]) Load() V     { return a.v.Load().(V) }
func (a *AtomicValue[V]) GetAndAdd(addFunc func(old V) V) V {
	for {
		old := a.v.Load().(V)
		if a.v.CompareAndSwap(old, addFunc(old)) {
			return old
		}
	}
}

func GetAndAdd(a *atomic.Int64, delta int64) int64 {
	for {
		old := a.Load()
		if a.CompareAndSwap(old, old+delta) {
			return old
		}
	}
}
