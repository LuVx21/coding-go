package types_x

import (
	"github.com/luvx21/coding-go/coding-common/maps_x"
)

type Slice[T any] []T
type NumberSlice[T Number] Slice[T]
type Stack[T any] []T

func (s *Stack[E]) IsEmpty() bool {
	return len(*s) == 0
}
func (s *Stack[E]) Push(e ...E) {
	*s = append(*s, e...)
}
func (s *Stack[E]) Peek() E {
	i := len(*s) - 1
	return (*s)[i]
}

func (s *Stack[E]) Pop() E {
	i := len(*s) - 1
	e := (*s)[i]
	*s = (*s)[:i]
	return e
}

type Map[K comparable, V any] map[K]V
type Set[E comparable] Map[E, struct{}]

func (m *Map[K, V]) Filter(f func(K, V) bool) Map[K, V] {
	return maps_x.Filter(*m, f)
}

func (m *Map[K, V]) Clone() Map[K, V] {
	return maps_x.Clone(*m)
}

func (m *Map[K, V]) Merge(source Map[K, V], replace bool) Map[K, V] {
	if *m == nil {
		*m = make(Map[K, V], len(source))
	}
	return maps_x.Merge(*m, source, replace)
}

func (s *Set[E]) Add(e ...E) {
	for _, _e := range e {
		(*s)[_e] = struct{}{}
	}
}

func (s *Set[E]) Remove(e ...E) {
	for _, _e := range e {
		delete(*s, _e)
	}
}

func (s *Set[E]) Contain(e E) bool {
	_, exist := (*s)[e]
	return exist
}

func (s *Slice[E]) Remove(e E) *Slice[E] {
	return s
}
