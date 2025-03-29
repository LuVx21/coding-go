package json

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jmespath-community/go-jmespath"
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
