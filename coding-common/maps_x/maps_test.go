package maps_x

import (
    "fmt"
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
