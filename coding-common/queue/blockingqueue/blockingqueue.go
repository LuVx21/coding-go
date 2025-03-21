package blockingqueue

import (
	"container/list"
	"sync"
)

type BlockingQueue[T any] struct {
	list *list.List
	lock sync.Mutex
	cond *sync.Cond
}

func New[T any]() *BlockingQueue[T] {
	q := &BlockingQueue[T]{list: list.New()}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *BlockingQueue[T]) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.list.Len() == 0
}

func (q *BlockingQueue[T]) Offer(items ...T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for _, item := range items {
		q.list.PushBack(item)
	}
	q.cond.Signal()
}

// Peek 返回队列头部的元素, 但不移除
func (q *BlockingQueue[T]) Peek() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.list.Len() == 0 {
		var zero T
		return zero, false
	}
	item := q.list.Front()
	return item.Value.(T), true
}

// Poll 从队列头部移除并返回元素, 如果队列为空则阻塞
func (q *BlockingQueue[T]) Poll() T {
	q.lock.Lock()
	defer q.lock.Unlock()

	for q.list.Len() == 0 {
		q.cond.Wait()
	}

	item := q.list.Front()
	q.list.Remove(item)
	return item.Value.(T)
}
