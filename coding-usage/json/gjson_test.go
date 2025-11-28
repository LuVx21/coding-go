package json

import (
	"fmt"
	"testing"

	"github.com/icloudza/fxjson"

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

func Test_fxjson_00(t *testing.T) {
	jsonData := []byte(`{
        "name": "Alice",
        "age": 30,
        "profile": {
            "city": "北京",
            "hobby": "coding"
        }
    }`)

	node := fxjson.FromBytes(jsonData)

	// 安全访问，使用默认值
	name := node.Get("name").StringOr("未知")
	age := node.Get("age").IntOr(0)
	city := node.GetPath("profile.city").StringOr("")

	fmt.Printf("姓名: %s, 年龄: %d, 城市: %s\n", name, age, city)
}
