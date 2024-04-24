package json

import (
    "encoding/json"
    "fmt"
    "github.com/jmespath-community/go-jmespath"
    "testing"
)

func Test_aa(t *testing.T) {
    _json := `
{
    "foo": {
        "bar": {
            "baz": [
                0,
                1,
                2,
                3,
                4
            ]
        }
    }
}
`
    var data interface{}
    _ = json.Unmarshal([]byte(_json), &data)
    search, _ := jmespath.Search("foo.bar.baz[2]", data)
    fmt.Println(search)
}
