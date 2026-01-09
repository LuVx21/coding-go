package json

import (
	"fmt"
	"testing"

	"github.com/icloudza/fxjson"

	"github.com/tidwall/gjson"
)

var jsonData = `{
		"name": "Alice",
		"age": 30,
		"profile": {
			"city": "北京",
			"hobby": "coding"
		},
		"aa": [
			{
				"id": 11
			},
			{
				"id": 12
			}
		],
		"bb": [1,2,3]
	}`

func Test_gjson_00(t *testing.T) {
	g := gjson.Parse(jsonData)
	s := g.Get("c")
	fmt.Println(s.Num, s.Raw, s.String())

	fmt.Println(g.Get("aa.#.id").Array())
}

func Test_fxjson_00(t *testing.T) {
	node := fxjson.FromString(jsonData)

	// 安全访问，使用默认值
	name := node.Get("name").StringOr("未知")
	age := node.Get("age").IntOr(0)
	city := node.GetPath("profile.city").StringOr("")
	node.Get("aa").ArrayForEach(func(index int, value fxjson.Node) bool {
		fmt.Println(value.Get("id").Int())
		return true
	})

	fmt.Printf("姓名: %s, 年龄: %d, 城市: %s\n", name, age, city)
}
