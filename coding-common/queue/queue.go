package queue

type Queue[T comparable] interface {
	IsEmpty() bool
	Offer(items ...T)
	Peek() (item T, ok bool)
	Poll() (item T, ok bool)
}
