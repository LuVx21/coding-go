package ds

type ListNode[T any] struct {
    Pre  *ListNode[T]
    Val  T
    Next *ListNode[T]
}

type TreeNode[T any] struct {
    Val   T
    Left  *TreeNode[T]
    Right *TreeNode[T]
}

type Node[T any] struct {
    Val      T
    Children []*Node[T]
}

func NewListNode[T any](v T, pre, next *ListNode[T]) *ListNode[T] {
    return &ListNode[T]{Pre: pre, Val: v, Next: next}
}
func NewTreeNode[T any](v T, left, right *TreeNode[T]) *TreeNode[T] {
    return &TreeNode[T]{Val: v, Left: left, Right: right}
}
func NewNode[T any](v T, nodes ...*Node[T]) *Node[T] {
    child := make([]*Node[T], len(nodes))
    copy(child, nodes)
    return &Node[T]{Val: v, Children: child}
}

func (m *TreeNode[T]) IsLeaf() bool {
    return m.Left == nil && m.Right == nil
}
