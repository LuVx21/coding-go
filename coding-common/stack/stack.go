package stack

type Stack[T any] interface {
	IsEmpty() bool
	Push(items ...T)
	Peek() (item T, ok bool)
	Pop() (item T, ok bool)
}
