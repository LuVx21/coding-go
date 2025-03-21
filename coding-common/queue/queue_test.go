package queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/queue/blockingqueue"
	. "github.com/luvx21/coding-go/coding-common/queue/listqueue"
	. "github.com/luvx21/coding-go/coding-common/queue/slicequeue"
)

func Test_array_queue(t *testing.T) {
	q := SliceQueue[string]{"foo", "bar"}
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

	q.Offer(3)
	q.Offer(2)
	q.Offer(1)

	for !q.IsEmpty() {
		item, _ := q.Poll()
		fmt.Println(item)
	}
}

func Test_block_queue(t *testing.T) {
	q := blockingqueue.New[int]()

	go func() {
		for i := 0; i < 5; i++ {
			q.Offer(i)
			fmt.Println("Produced:", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			item := q.Poll()
			fmt.Println("Consumed:", item)
		}
	}()

	time.Sleep(5 * time.Second)
}
