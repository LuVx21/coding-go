package slicequeue

import (
	"fmt"
	"strings"
)

type Queue[T any] []T

func (s *Queue[E]) IsEmpty() bool {
	return len(*s) == 0
}
func (s *Queue[E]) Offer(e ...E) {
	*s = append(*s, e...)
}
func (s *Queue[E]) Peek() E {
	return (*s)[0]
}
func (s *Queue[E]) Poll() E {
	e := (*s)[0]
	*s = (*s)[1:]
	return e
}
func (s *Queue[E]) OfferFirst(e ...E) {
	*s = append(e, *s...)
}
func (s *Queue[E]) PeekLast() E {
	i := len(*s) - 1
	return (*s)[i]
}
func (s *Queue[E]) PollLast() E {
	i := len(*s) - 1
	e := (*s)[i]
	*s = (*s)[:i]
	return e
}

func (s *Queue[E]) String() string {
	var sb strings.Builder
	for i, item := range *s {
		sb.WriteString(fmt.Sprintf("%v", item))
		if i < len(*s)-1 {
			sb.WriteString("←")
		}
	}
	return sb.String()
}
