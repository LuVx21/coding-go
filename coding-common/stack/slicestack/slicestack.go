package slicestack

import (
	"fmt"
	"slices"
	"strings"
)

type SliceStack[T any] []T

func (s *SliceStack[E]) IsEmpty() bool {
	return len(*s) == 0
}
func (s *SliceStack[E]) Push(e ...E) {
	*s = append(*s, e...)
}
func (s *SliceStack[E]) Peek() E {
	i := len(*s) - 1
	return (*s)[i]
}

func (s *SliceStack[E]) Pop() E {
	i := len(*s) - 1
	e := (*s)[i]
	*s = (*s)[:i]
	return e
}

func (s *SliceStack[E]) String() string {
	var sb strings.Builder
	for i := range slices.Backward(*s) {
		fmt.Fprintf(&sb, "%v\n", (*s)[i])
		if i > 0 {
			sb.WriteString("↑\n")
		}
	}
	return sb.String()
}
