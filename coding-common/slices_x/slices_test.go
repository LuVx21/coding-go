package slices_x

import (
    "fmt"
    "testing"
)

func Test_Partition(t *testing.T) {
    ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
    r := Partition(ints, 5)
    fmt.Println(r)
}

func Test_01(t *testing.T) {

    bs := []bool{true, false, true}
    fmt.Println(ClearZero(bs))

    is := []int{1, 2, 0, 0, 3}
    fmt.Println(ClearZero(is))

    strs := []string{"", "foo", "bar", "baz"}
    fmt.Println(ClearZero(strs))

    ms := []map[string]string{
        {"foo": "foo"},
        nil,
        {"bar": "bar"},
    }
    fmt.Println(ClearZeroRef(ms))
}
