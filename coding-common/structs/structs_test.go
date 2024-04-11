package structs

import (
    "fmt"
    "testing"
)

type UserInfo struct {
    Name     string   `json:"name"`
    Age      int      `json:"age"`
    Age1     *int     `json:"age1"`
    Profile  Profile  `json:"profile"`
    Profile1 *Profile `json:"profile1"`
}
type Profile struct {
    Hobby string `json:"hobby"`
}

var a = 18
var u1 = UserInfo{
    Name:     "foo-bar",
    Age:      a,
    Age1:     &a,
    Profile:  Profile{"双色球"},
    Profile1: &Profile{"双色球"},
}

func Test_01(t *testing.T) {
    m1, _ := ToMap(&u1, "json")
    for k, v := range m1 {
        fmt.Printf("k: %v v: %v 类型: %T\n", k, v, v)
    }

    m2, _ := ToSingleMap(&u1, "json")
    for k, v := range m2 {
        fmt.Printf("k: %v v: %v 类型: %T\n", k, v, v)
    }
}
