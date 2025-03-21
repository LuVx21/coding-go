package set

import "github.com/luvx21/coding-go/coding-common/maps_x"

type Set[E comparable] maps_x.Map[E, struct{}]

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
