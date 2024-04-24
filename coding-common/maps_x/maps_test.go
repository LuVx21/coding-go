package maps_x

import (
    "fmt"
    . "github.com/luvx21/coding-go/coding-common/common_x/types_x"
    "reflect"
    "testing"
)

func Test_Join(t *testing.T) {
    m := map[string]interface{}{
        "a": "aa",
        "b": "bb",
    }
    join := Join(m, "=", ";")
    fmt.Println(join)
}

func Test_RemoveIf(t *testing.T) {
    m := map[string]interface{}{
        "a": "aa",
        "b": "bb",
    }
    RemoveIf(m, func(k string, v interface{}) bool {
        return k == "a"
    })
    fmt.Println(m)
}

func Test_01(t *testing.T) {
    m := make(map[string]interface{})
    m["b"] = "b"

    if b := GetOrDefault(m, "b", "haha"); true {
        fmt.Println(b)
    }
    if c := GetOrDefault(m, "c", "haha"); true {
        fmt.Println(c)
    }
}

func Test_02(t *testing.T) {
    m := make(map[string]interface{})

    m["a"] = nil
    m["b"] = "bb"

    absent := ComputeIfAbsent(m, "a", func(s string) interface{} {
        return s + "-100"
    })
    fmt.Println(absent)

    m["c"] = nil
    m["d"] = "dd"
    present := ComputeIfPresent(m, "d", func(s string, i interface{}) interface{} {
        return fmt.Sprintf("%s-%s-200", s, i)
    })
    fmt.Println(present)

    m["e"] = "ee"
    compute := Compute(m, "f", func(s string, i interface{}) interface{} {
        return fmt.Sprintf("%s-%s-200", s, i)
    })
    fmt.Println(compute)
}

func TestGetByKey(t *testing.T) {
    var i int64 = 999
    m := Map[string]{"a": i, "m": Map[string]{"mm": "mm"}, "s": []string{"ss"}}

    v, err := GetByKey(m, "a", reflect.Int64, 0)
    fmt.Println(v, err)

    v, err = GetInt64(m, "a", -1)
    fmt.Println(v, err)

    v, err = GetMap(m, "mm", map[any]any{"test": "test"})
    fmt.Println(v, err)

    v, err = GetSlice(m, "s", []string{"sss"})
    fmt.Println(v, err)
}
