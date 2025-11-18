package syncx

import "sync"

type SyncMap[K comparable, V any] struct {
	m *sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	s := new(SyncMap[K, V])
	s.m = new(sync.Map)
	return s
}

func (s *SyncMap[K, V]) Store(key K, val V) {
	s.m.Store(key, val)
}
func (s *SyncMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}
func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	if v, ok := s.m.Load(key); ok {
		return v.(V), ok
	}
	var zero V
	return zero, false
}
func (s *SyncMap[K, V]) Range(f func(key K, val V) bool) {
	f2 := func(key, value any) bool { return f(key.(K), value.(V)) }
	s.m.Range(f2)
}
