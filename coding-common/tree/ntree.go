package tree

type Node[T any] struct {
	Val      T
	Children []*Node[T]
}

func NewNode[T any](v T, nodes ...*Node[T]) *Node[T] {
	child := make([]*Node[T], len(nodes))
	copy(child, nodes)
	return &Node[T]{Val: v, Children: child}
}
