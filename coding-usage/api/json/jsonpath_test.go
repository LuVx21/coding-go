package json

import (
    "encoding/json"
    "fmt"
    "testing"
    "github.com/jmespath-community/go-jmespath"
)

func Test_a(t *testing.T) {
    var jsondata = []byte(`{"foo": {"bar": {"baz": [0, 1, 2, 3, 4]}}}`) // your data
    var data interface{}
    err := json.Unmarshal(jsondata, &data)
    result, err := jmespath.Search("foo.bar.baz[2]", data)
}
