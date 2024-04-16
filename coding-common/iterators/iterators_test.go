package iterators

import (
    "github.com/luvx21/coding-go/coding-common/logs"
    "strconv"
    "testing"
)

type Item struct {
    id   int
    name string
}

func dao1(cursor int, limit int) []Item {
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
    logs.Log.Printf("cursor:%d limit:%d data:%v", cursor, limit, res)
    return res
}

func Test_01(t *testing.T) {
    var limit = 10
    iterator := NewCursorIterator(
        0,
        false,
        func(id int) []Item {
            return dao1(id, limit)
        },
        func(items []Item) int {
            if len(items) < limit {
                return -1
            }
            return items[len(items)-1].id + 1
        },
        func(i int) bool {
            return i < 0 || i > 47
        },
    )

    iterator.forEachRemaining(func(item Item) {
        logs.Log.Printf("next:%v", item)
    })
}
