package dao

import (
    "fmt"
    "luvx/gin/dao/common_kv"
    "testing"
)

func Test_01(t *testing.T) {
    m := common_kv.Get(7, "foo", "bar")
    for _, e := range m {
        fmt.Println(e.CommonKey)
    }
}
