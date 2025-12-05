package ds

type ListNode[T any] struct {
	Val       T
	Pre, Next *ListNode[T]
}

func NewListNode[T any](v T, pre, next *ListNode[T]) *ListNode[T] {
	return &ListNode[T]{Pre: pre, Val: v, Next: next}
}
