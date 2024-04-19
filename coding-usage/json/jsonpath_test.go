package json

import (
    "encoding/json"
    "github.com/jmespath-community/go-jmespath"
    "testing"
)

func Test_aa(t *testing.T) {
    var jsondata = []byte(`{"foo": {"bar": {"baz": [0, 1, 2, 3, 4]}}}`) // your data
    var data interface{}
    _ = json.Unmarshal(jsondata, &data)
    _, _ = jmespath.Search("foo.bar.baz[2]", data)
}
