package listqueue

import (
	"container/list"
)

type ListQueue[T any] struct {
	items *list.List
}

func NewListQueue[T any]() *ListQueue[T] {
	return &ListQueue[T]{items: list.New()}
}

func (q *ListQueue[T]) Enqueue(item T) {
	q.items.PushBack(item)
}

func (q *ListQueue[T]) Dequeue() (T, bool) {
	if q.IsEmpty() {
		var r T
		return r, false
	}
	item := q.items.Front()
	q.items.Remove(item)
	return item.Value.(T), true
}

func (q *ListQueue[T]) IsEmpty() bool {
	return q.items.Len() == 0
}
