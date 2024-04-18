package jsons

import (
    "fmt"
    "testing"
)

func Test_01(t *testing.T) {
    s := "{\"a\": true}"
    toMap, err := JsonStringToMap[string, bool, map[string]bool](s)
    fmt.Println(toMap, err)
    s1 := "[\"a\", \"b\"]"
    array, err1 := JsonStringToArray[string, []string](s1)
    fmt.Println(array, err1)
}
