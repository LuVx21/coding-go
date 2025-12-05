package sync_sets

import (
	"fmt"
	"strings"
	"sync"
)

type set[E comparable] struct {
	m sync.Map
}

func NewSet[E comparable](e ...E) *set[E] {
	set := &set[E]{}
	set.Add(e...)
	return set
}

func (s *set[E]) Len() int {
	count := 0
	s.m.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func (s *set[E]) Add(e ...E) {
	for _, _e := range e {
		s.m.Store(_e, struct{}{})
	}
}

func (s *set[E]) Remove(e ...E) {
	for _, _e := range e {
		s.m.Delete(_e)
	}
}

func (s *set[E]) Contains(e E) bool {
	_, exist := s.m.Load(e)
	return exist
}

func (s *set[E]) ToSlice() []E {
	r := make([]E, 0)
	s.m.Range(func(key, _ any) bool {
		r = append(r, key.(E))
		return true
	})
	return r
}

func (s *set[E]) String() string {
	var sb strings.Builder
	s.m.Range(func(key, _ any) bool {
		sb.WriteString(fmt.Sprintf("%v\n", key))
		return true
	})
	return sb.String()
}
