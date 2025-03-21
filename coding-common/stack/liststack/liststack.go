package liststack

import (
	"container/list"
	"fmt"
	"strings"
)

type ListStack[T any] struct {
	items *list.List
}

func NewListStack[T any](items ...T) *ListStack[T] {
	l := &ListStack[T]{items: list.New()}
	l.Push(items...)
	return l
}

func (s *ListStack[E]) IsEmpty() bool {
	return s.items.Len() == 0
}
func (s *ListStack[E]) Push(e ...E) {
	for _, item := range e {
		s.items.PushBack(item)
	}
}
func (s *ListStack[E]) Peek() (E, bool) {
	if s.IsEmpty() {
		var r E
		return r, false
	}
	item := s.items.Back()
	return item.Value.(E), true
}
func (s *ListStack[E]) Pop() (E, bool) {
	if s.IsEmpty() {
		var r E
		return r, false
	}
	item := s.items.Back()
	s.items.Remove(item)
	return item.Value.(E), true
}

func (s *ListStack[E]) String() string {
	var sb strings.Builder
	for cur := s.items.Back(); cur != nil; cur = cur.Prev() {
		sb.WriteString(fmt.Sprintf("%v\n", cur.Value))
		if cur.Prev() != nil {
			sb.WriteString("↑\n")
		}
	}
	return sb.String()
}
