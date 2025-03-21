package types_x

import (
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
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

func Test_Set(t *testing.T) {
	s := Set[string]{"foobar": struct{}{}}
	s.Add("foo", "bar")
	fmt.Println(s.Contain("foo"), s.Contain("bar1"))
}
