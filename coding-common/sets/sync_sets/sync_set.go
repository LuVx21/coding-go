package sync_sets

import (
	"fmt"
	"strings"
	"sync"
)

type Set[E comparable] struct {
	m sync.Map
}

func NewSet[E comparable](e ...E) *Set[E] {
	set := &Set[E]{}
	set.Add(e...)
	return set
}

func (s *Set[E]) Len() int {
	count := 0
	s.m.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func (s *Set[E]) Add(e ...E) {
	for _, _e := range e {
		s.m.Store(_e, struct{}{})
	}
}

func (s *Set[E]) Remove(e ...E) {
	for _, _e := range e {
		s.m.Delete(_e)
	}
}

func (s *Set[E]) Clear() {
	s.m.Range(func(key, value any) bool {
		s.m.Delete(key)
		return true
	})
}

func (s *Set[E]) Contains(e E) bool {
	_, exist := s.m.Load(e)
	return exist
}

func (s *Set[E]) ToSlice() []E {
	r := make([]E, 0)
	s.m.Range(func(key, _ any) bool {
		r = append(r, key.(E))
		return true
	})
	return r
}

func (s *Set[E]) String() string {
	var sb strings.Builder
	s.m.Range(func(key, _ any) bool {
		fmt.Fprintf(&sb, "%v\n", key)
		return true
	})
	return sb.String()
}
