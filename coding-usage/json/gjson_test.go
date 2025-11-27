package json

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

func Test_gjson_00(t *testing.T) {
	_json := `
{
    "a": {
        "b": "bbb"
    },
	"b": [1,2,3],
	"c": 101
}
`
	s := gjson.Get(_json, "c")
	fmt.Println(s.Num, s.Raw, s.String())
}
