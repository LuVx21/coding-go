package listqueue

import (
	"container/list"
	"fmt"
	"strings"
)

type ListQueue[T any] struct {
	items *list.List
}

func NewListQueue[T any](items ...T) *ListQueue[T] {
	l := &ListQueue[T]{items: list.New()}
	l.Offer(items...)
	return l
}

func (q *ListQueue[T]) IsEmpty() bool {
	return q.items.Len() == 0
}

// Offer 队尾添加
func (q *ListQueue[T]) Offer(items ...T) {
	for _, item := range items {
		q.items.PushBack(item)
	}
}

// Peek 查看队首
func (q *ListQueue[T]) Peek() (T, bool) {
	if q.IsEmpty() {
		var r T
		return r, false
	}
	item := q.items.Front()
	return item.Value.(T), true
}

// Poll 查看移除队首
func (q *ListQueue[T]) Poll() (T, bool) {
	if q.IsEmpty() {
		var r T
		return r, false
	}
	item := q.items.Front()
	q.items.Remove(item)
	return item.Value.(T), true
}

// OfferFirst 队首添加
func (q *ListQueue[T]) OfferFirst(items ...T) {
	for _, item := range items {
		q.items.PushFront(item)
	}
}

// PeekLast 查看队尾
func (q *ListQueue[T]) PeekLast() (T, bool) {
	if q.IsEmpty() {
		var r T
		return r, false
	}
	item := q.items.Back()
	return item.Value.(T), true
}

// PollLast 查看移除队尾
func (q *ListQueue[T]) PollLast() (T, bool) {
	if q.IsEmpty() {
		var r T
		return r, false
	}
	item := q.items.Back()
	q.items.Remove(item)
	return item.Value.(T), true
}

func (q *ListQueue[E]) String() string {
	var sb strings.Builder
	for cur := q.items.Front(); cur != nil; cur = cur.Next() {
		sb.WriteString(fmt.Sprintf("%v", cur.Value))
		if cur.Next() != nil {
			sb.WriteString("←")
		}
	}
	return sb.String()
}
