package iterators

import (
	"fmt"
	"github.com/luvx21/coding-go/coding-common/common_x/pairs"
	"log"
	"strconv"
	"testing"
)

type Item struct {
	id   int
	name string
}

func dao0(cursor int, limit int) pairs.Pair[[]Item, int] {
	res := make([]Item, 0)
	var cnt = 1
	for i := 0; i <= 112 && cnt <= limit; i++ {
		if i >= cursor {
			item := Item{
				id:   i,
				name: "No." + strconv.Itoa(i),
			}
			res = append(res, item)
			cnt++
		}
	}
	log.Printf("cursor:%d limit:%d data:%v", cursor, limit, res)
	return pairs.NewPair(res, cursor+limit)
}

func dao1(cursor int, limit int) []Item {
	return dao0(cursor, limit).K
}

func dao2(pageNo int, limit int) pairs.Pair[[]Item, int] {
	return dao0((pageNo-1)*limit, limit)
}

func dao3(pageNo int, limit int) []Item {
	return dao2(pageNo, limit).K
}

func Test_00(t *testing.T) {
	const limit = 10
	iterator := NewCursorIterator[Item, int, pairs.Pair[[]Item, int]](
		0,
		false,
		func(id int) pairs.Pair[[]Item, int] {
			return dao0(id, limit)
		},
		func(curId int, p pairs.Pair[[]Item, int]) int {
			items := p.K
			if len(items) < limit {
				return -1
			}
			return items[len(items)-1].id + 1
		},
		func(p pairs.Pair[[]Item, int]) []Item {
			return p.K
		},
		func(i int) bool {
			return i < 0 || i > 47
		},
	)

	iterator.ForEachRemaining(func(item Item) {
		log.Printf("next:%v", item)
	})
}

func Test_01(t *testing.T) {
	const limit = 10
	iterator := NewCursorIteratorSimple[Item, int](
		0,
		false,
		func(id int) []Item {
			return dao1(id, limit)
		},
		func(curId int, items []Item) int {
			fmt.Println("本次结果:", curId, items)
			if len(items) < limit {
				return -1
			}
			return items[len(items)-1].id + 1
		},
		func(i int) bool {
			return i < 0 || i > 47
		},
	)

	iterator.ForEachRemaining(func(item Item) {
		log.Printf("next:%v", item)
	})
}

func Test_page(t *testing.T) {
	const limit = 10
	iterator :=
		NewPageIterator[Item, pairs.Pair[[]Item, int]](
			0,
			false,
			func(pageNo int) pairs.Pair[[]Item, int] {
				return dao2(pageNo, limit)
			},
			func(p pairs.Pair[[]Item, int]) []Item {
				return p.K
			},
			func(pageNo int) bool {
				return pageNo < 0 || pageNo > 5
			},
		)

	// NewPageIteratorSimple[Item](
	//    0,
	//    false,
	//    func(pageNo int) []Item {
	//        return dao3(pageNo, limit)
	//    },
	//    func(pageNo int) bool {
	//        return pageNo < 0 || pageNo > 5
	//    },
	// )

	iterator.ForEachRemaining(func(item Item) {
		log.Printf("next:%v", item)
	})
}
