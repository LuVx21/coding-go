package json

import (
	"encoding/json/v2"
	"fmt"
	"testing"

	"github.com/jmespath-community/go-jmespath"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
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
	var data any
	_ = json.Unmarshal([]byte(_json), &data)
	search, _ := jmespath.Search("foo.bar.baz[2]", data)
	fmt.Println(search)
}

func Test_ojg_00(t *testing.T) {
	obj, _ := oj.ParseString(`{
        "a":[
            {"x":1,"y":2,"z":3},
            {"x":2,"y":4,"z":6}
        ]
    }`)

	x, _ := jp.ParseString("a[?(@.x > 1)].y")
	ys := x.Get(obj)
	fmt.Println(ys)
}
