package t

import "iter"

type (
	// ListNode 链表节点
	ListNode[T any] struct {
		Val       T
		Pre, Next *ListNode[T]
	}
	List[T any] struct {
		Head, Tail *ListNode[T]
	}
)

// ------------------------------------------------------------------------------------------------------------------------
func (lst *List[T]) Push(vs ...T) {
	if len(vs) == 0 {
		return
	}
	tem := vs
	if lst.Tail == nil {
		lst.Head = &ListNode[T]{Val: vs[0]}
		lst.Tail = lst.Head
		tem = vs[1:]
	}
	for _, v := range tem {
		lst.Tail.Next = &ListNode[T]{Val: v}
		lst.Tail = lst.Tail.Next
	}
}
func (lst *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := lst.Head; e != nil; e = e.Next {
			if !yield(e.Val) {
				return
			}
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------

func NewListNode[T any](v T, pre, Next *ListNode[T]) *ListNode[T] {
	return &ListNode[T]{Pre: pre, Val: v, Next: Next}
}
func (n *ListNode[T]) Data() T             { return n.Val }
func (n *ListNode[T]) Prev() *ListNode[T]  { return n.Pre }
func (n *ListNode[T]) NextN() *ListNode[T] { return n.Next }
