package sortx

import (
    "sort"
    "testing"
    "time"
)

type item struct {
    t time.Time
}

type items []item

func (s items) Len() int {
    return len(s)
}
func (s items) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s items) Less(i, j int) bool {
    return s[i].t.After(s[j].t)
}

func beforeAfter(s []item) func() {
    for _, t := range s {
        println(t.t.Format("2006-01-02 15:04:05.999999"))
    }

    return func() {
        println("--------------------")
        for _, t := range s {
            println(t.t.Format("2006-01-02 15:04:05.999999"))
        }
    }
}

func Test_sort_00(t *testing.T) {
    mySlice := []item{{time.Now()}, {time.Now().Add(time.Hour)}, {time.Now().Add(-time.Hour)}}
    defer beforeAfter(mySlice)()

    // 不需要实现方式
    sort.Slice(mySlice, func(i, j int) bool {
        return mySlice[i].t.After(mySlice[j].t)
    })
}

func Test_sort_01(t *testing.T) {
    mySlice := items{{time.Now()}, {time.Now().Add(time.Hour)}, {time.Now().Add(-time.Hour)}}
    defer beforeAfter(mySlice)()

    // 需要实现接口
    sort.Sort(mySlice)
}

func Test_sort_02(t *testing.T) {
    mySlice := []item{{time.Now()}, {time.Now().Add(time.Hour)}, {time.Now().Add(-time.Hour)}}
    defer beforeAfter(mySlice)()

    wrapper := SortWrapper[item]{
        items: mySlice,
        By:    func(l, r *item) bool { return l.t.After(r.t) },
    }
    // 需要实现接口, 自定义排序方式
    sort.Sort(&wrapper)
}

func Test_sort_03(t *testing.T) {
    mySlice := []item{{time.Now()}, {time.Now().Add(time.Hour)}, {time.Now().Add(-time.Hour)}}
    defer beforeAfter(mySlice)()

    Sort(mySlice, func(l, r *item) bool {
        return l.t.After(r.t)
    })
}
