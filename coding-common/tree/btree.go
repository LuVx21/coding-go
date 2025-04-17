package tree

import (
	"fmt"
	"strings"
)

type TreeNode[T any] struct {
	Val         T
	Left, Right *TreeNode[T]
}

func NewTreeNode[T any](v T, left, right *TreeNode[T]) *TreeNode[T] {
	return &TreeNode[T]{Val: v, Left: left, Right: right}
}

func NewCBT[T any](arr ...T) *TreeNode[T] {
	_len := len(arr)
	if _len == 0 {
		return nil
	} else if _len == 1 {
		return &TreeNode[T]{Val: arr[0]}
	}

	list := make([]*TreeNode[T], _len)
	for i, v := range arr {
		list[i] = &TreeNode[T]{Val: v}
	}
	for i := 0; 2*i+1 < _len; i++ {
		n := list[i]
		if n == nil {
			continue
		}
		n.Left = list[2*i+1]
		if 2*i+2 < _len {
			n.Right = list[2*i+2]
		}
	}
	return list[0]
}

func (n *TreeNode[T]) String() string {
	var sb strings.Builder
	n.printTree(&sb, "", "")
	return sb.String()
}

func (n *TreeNode[T]) printTree(sb *strings.Builder, prefix string, childPrefix string) {
	if n == nil {
		return
	}
	sb.WriteString(prefix)
	sb.WriteString(fmt.Sprintf("%v\n", n.Val))

	if n.Left != nil || n.Right != nil {
		n.Left.printTree(sb, childPrefix+"├", childPrefix+"│")
		n.Right.printTree(sb, childPrefix+"└", childPrefix+"  ")
	}
}

func (m *TreeNode[T]) IsLeaf() bool {
	return m.Left == nil && m.Right == nil
}

func (root *TreeNode[T]) String1() string {
	if root == nil {
		return ""
	}

	levelRows := make([][]T, 0)
	queue := make([]*TreeNode[T], 0)

	for k := 0; len(queue) > 0; k++ {
		var isBreak bool = true
		size := len(queue)
		row := make([]T, 0)
		for range size {
			node := queue[len(queue)-1]
			queue = queue[:len(queue)-1]
			if node == nil {
				// TODO
				// row = append(row, math.MinInt)
				queue = append(queue, nil, nil)
			} else {
				isBreak = isBreak && node.IsLeaf()
				row = append(row, node.Val)
				queue = append(queue, node.Left, node.Right)
			}
		}
		levelRows = append(levelRows, row)
		if isBreak {
			break
		}
	}
	line := strings.Repeat("-", 150)
	sb := dumpTreeFormat0(levelRows)
	sb.WriteString(line)
	sb.WriteString("\n")
	sb.WriteString(dumpTreeFormat(levelRows))
	sb.WriteString(line)
	sb.WriteString("\n")
	return sb.String()
}

func dumpTreeFormat0[T any](rows [][]T) strings.Builder {
	var sb strings.Builder
	return sb
}
func dumpTreeFormat[T any](rows [][]T) string {
	return ""
}
