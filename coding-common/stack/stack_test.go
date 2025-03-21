package stack

import (
	"fmt"
	"testing"

	. "github.com/luvx21/coding-go/coding-common/stack/liststack"
	. "github.com/luvx21/coding-go/coding-common/stack/slicestack"
)

func Test_slice_stack(t *testing.T) {
	s := Stack[string]{"foo", "bar"}
	fmt.Println(s)
	s.Push("aaa")
	fmt.Println(s, s.Peek(), "----------")
	fmt.Println(s)
	top := s.Pop()
	fmt.Println(s, top, s.IsEmpty(), "----------")
	s.Pop()
	s.Pop()
	fmt.Println(s, s.IsEmpty())
}

func Test_list_stack(t *testing.T) {
	s := NewListStack[string]("foo", "bar")
	s.Push("aaa")
	fmt.Println(s)
}
