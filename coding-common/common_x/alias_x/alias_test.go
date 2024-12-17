package alias_x

import (
    "fmt"
    "testing"
)

func Test_alias_00(t *testing.T) {
    a := Slice[int]{1, 2, 3}
    b := Slice[any]{1, 2, 3, "a"}

    c := SliceAny{1, 2, 3, "a"}
    fmt.Println(a, b, c)

    m := MapComparable2Any[uint]{1: "a", 2: "b"}
    fmt.Println(m)
}
