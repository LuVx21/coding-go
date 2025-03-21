package blockingstack

import (
	"container/list"
	"sync"
)

type BlockingStack[T any] struct {
	list *list.List
	lock sync.Mutex
	cond *sync.Cond
}

func New[T any]() *BlockingStack[T] {
	s := &BlockingStack[T]{list: list.New()}
	s.cond = sync.NewCond(&s.lock)
	return s
}

func (s *BlockingStack[T]) IsEmpty() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.list.Len() == 0
}

func (s *BlockingStack[T]) Push(items ...T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range items {
		s.list.PushBack(item)
	}
	s.cond.Signal()
}

// Peek 返回栈顶的元素，但不移除
func (s *BlockingStack[T]) Peek() (T, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.list.Len() == 0 {
		var zero T
		return zero, false
	}
	item := s.list.Back()
	return item.Value.(T), true
}

func (s *BlockingStack[T]) Pop() (T, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for s.list.Len() == 0 {
		s.cond.Wait()
	}

	item := s.list.Back()
	s.list.Remove(item)
	return item.Value.(T), true
}
