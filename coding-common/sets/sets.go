package sets

import (
	"fmt"
	"strings"
)

type Set[E comparable] map[E]struct{}

func NewSet[E comparable](e ...E) *Set[E] {
	set := &Set[E]{}
	set.Add(e...)
	return set
}

func (s *Set[E]) Len(e ...E) int {
	return len(*s)
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

func (s *Set[E]) ToSlice(e E) []E {
	r := make([]E, s.Len())
	for k := range *s {
		r = append(r, k)
	}
	return r
}

func (s *Set[E]) String() string {
	var sb strings.Builder
	for k := range *s {
		sb.WriteString(fmt.Sprintf("%v\n", k))
	}
	return sb.String()
}
