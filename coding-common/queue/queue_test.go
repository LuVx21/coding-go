package queue

import (
	"fmt"
	"testing"

	. "github.com/luvx21/coding-go/coding-common/queue/listqueue"
	. "github.com/luvx21/coding-go/coding-common/queue/slicequeue"
)

func Test_array_queue(t *testing.T) {
	q := Queue[string]{"foo", "bar"}
	fmt.Println(q)
	q.Offer("aaa")
	fmt.Println(q, q.Peek(), "----------")
	fmt.Println(q)
	top := q.Poll()
	fmt.Println(q, top, q.IsEmpty(), "----------")
	q.Poll()
	q.Poll()
	fmt.Println(q, q.IsEmpty())
}

func Test_list_queue(t *testing.T) {
	q := NewListQueue[int]()

	q.Enqueue(3)
	q.Enqueue(2)
	q.Enqueue(1)

	for !q.IsEmpty() {
		item, _ := q.Dequeue()
		fmt.Println(item)
	}
}
