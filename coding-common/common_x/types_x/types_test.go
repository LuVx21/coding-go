package types_x

import (
    "fmt"
    "github.com/bytedance/sonic"
    "testing"
)

func Test_00(t *testing.T) {
    _json := `
{
  "foo": "bar",
  "bar": 2
}
`
    m := Map[string, any]{}
    _ = sonic.Unmarshal([]byte(_json), &m)
    fmt.Println(m)

    merge := m.Merge(Map[string, any]{"aaa": "bbb"}, true)
    fmt.Println(merge)

    nm := merge.Filter(func(k string, v any) bool {
        return k == "aaa"
    })
    fmt.Println(nm)
}

func Test_Stack(t *testing.T) {
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

func Test_Queue(t *testing.T) {
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

func Test_Set(t *testing.T) {
    s := Set[string]{"foobar": struct{}{}}
    s.Add("foo", "bar")
    fmt.Println(s.Contain("foo"), s.Contain("bar1"))
}
