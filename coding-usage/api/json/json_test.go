package json

import (
    "encoding/json"
    "fmt"
    "testing"
)

func Test_a(t *testing.T) {
    //a()
    //b()
    c()
}

func Test_b(t *testing.T) {
    body := `{
    "data": {
        "body": "aaaa"
    },
    "aa": 1
}`
    fmt.Println(body)
    m := make(map[string]interface{})
    json.Unmarshal([]byte(body), &m)

    a2 := m["data"]
    if kvs, ok := a2.(map[string]interface{}); ok {
        fmt.Println(a2, kvs["body"])
    }
    //fmt.Println(a2, reflect.TypeOf(a2), reflect.TypeOf(m), a2["body"])
    //a2["body"]
}
