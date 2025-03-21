package slicestack

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
