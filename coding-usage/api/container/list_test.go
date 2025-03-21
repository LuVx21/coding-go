package container

import (
	"container/list"
	"fmt"
	"testing"

	listx "github.com/smallnest/exp/container/list"
)

func Test_list_00(t *testing.T) {
	l := list.New()
	for i := 0; i <= 16; i++ {
		l.PushBack(i + 1)
	}

	for cur := l.Front(); cur != nil; cur = cur.Next() {
		fmt.Println(cur.Value)
	}
}

func Test_list_01(t *testing.T) {
	l := listx.New[string]()
	l.PushBack("bb")
	front := l.Front()
	fmt.Println(front.Value)
}
