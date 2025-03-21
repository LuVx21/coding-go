package stack

import (
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/stack/blockingstack"
	. "github.com/luvx21/coding-go/coding-common/stack/liststack"
	. "github.com/luvx21/coding-go/coding-common/stack/slicestack"
)

func Test_slice_stack(t *testing.T) {
	s := SliceStack[string]{"foo", "bar"}
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

func Test_block_stack(t *testing.T) {
	s := blockingstack.New[int]()

	// 生产者
	go func() {
		for i := 0; i < 5; i++ {
			s.Push(i)
			fmt.Println("Pushed:", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// 消费者
	go func() {
		for i := 0; i < 5; i++ {
			item, _ := s.Pop()
			fmt.Println("Popped:", item)
		}
	}()

	time.Sleep(5 * time.Second)
}
