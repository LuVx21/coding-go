package sets

import (
	"fmt"
	"strings"
)

type set[E comparable] map[E]struct{}

func NewSet[E comparable](e ...E) *set[E] {
	set := &set[E]{}
	set.Add(e...)
	return set
}

func (s *set[E]) Len() int {
	return len(*s)
}

func (s *set[E]) Add(e ...E) {
	for _, _e := range e {
		(*s)[_e] = struct{}{}
	}
}

func (s *set[E]) Remove(e ...E) {
	for _, _e := range e {
		delete(*s, _e)
	}
}

func (s *set[E]) Contain(e E) bool {
	_, exist := (*s)[e]
	return exist
}

func (s *set[E]) ToSlice(e E) []E {
	r := make([]E, s.Len())
	for k := range *s {
		r = append(r, k)
	}
	return r
}

func (s *set[E]) String() string {
	var sb strings.Builder
	for k := range *s {
		sb.WriteString(fmt.Sprintf("%v\n", k))
	}
	return sb.String()
}
