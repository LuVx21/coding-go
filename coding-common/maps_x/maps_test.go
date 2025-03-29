package maps_x

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	// json "github.com/bytedance/sonic"
)

func Test_map(t *testing.T) {
	_json := `
{
  "foo": "bar",
  "bar": 2
}
`
	m := Map[string, any]{}
	_ = json.Unmarshal([]byte(_json), &m)
	fmt.Println(m)

	merge := m.Merge(Map[string, any]{"aaa": "bbb"}, true)
	fmt.Println(merge)

	nm := merge.Filter(func(k string, v any) bool {
		return k == "aaa"
	})
	fmt.Println(nm)
}

func Test_Join(t *testing.T) {
	m := map[string]any{
		"a": "aa",
		"b": "bb",
	}
	join := Join(m, "=", ";")
	fmt.Println(join)
}

func Test_RemoveIf(t *testing.T) {
	m := map[string]any{
		"a": "aa",
		"b": "bb",
	}
	RemoveIf(m, func(k string, v any) bool {
		return k == "a"
	})
	fmt.Println(m)
}

func Test_01(t *testing.T) {
	m := make(map[string]any)
	m["b"] = "b"

	if b := GetOrDefault(m, "b", "haha"); true {
		fmt.Println(b)
	}
	if c := GetOrDefault(m, "c", "haha"); true {
		fmt.Println(c)
	}
}

func Test_02(t *testing.T) {
	m := make(map[string]any)

	m["a"] = nil
	m["b"] = "bb"

	absent := ComputeIfAbsent(m, "a", func(s string) any {
		return s + "-100"
	})
	fmt.Println(absent)

	m["c"] = nil
	m["d"] = "dd"
	present := ComputeIfPresent(m, "d", func(s string, i any) any {
		return fmt.Sprintf("%s-%s-200", s, i)
	})
	fmt.Println(present)

	m["e"] = "ee"
	compute := Compute(m, "f", func(s string, i any) any {
		return fmt.Sprintf("%s-%s-200", s, i)
	})
	fmt.Println(compute)
}

func TestGetByKey(t *testing.T) {
	var i int64 = 999
	m := map[string]any{"a": i, "m": map[string]any{"mm": "mm"}, "s": []string{"ss"}}

	v, err := GetByKey(m, "a", reflect.Int64, 0)
	fmt.Println(v, err)

	v, err = GetInt64(m, "a", -1)
	fmt.Println(v, err)

	v, err = GetMap(m, "mm", map[any]any{"test": "test"})
	fmt.Println(v, err)

	v, err = GetSlice(m, "s", []string{"sss"})
	fmt.Println(v, err)
}

func Test_empty_01(t *testing.T) {
	e := Empty[string, string]()
	fmt.Println(e)
}
