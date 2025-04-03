package json

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
	goJson "github.com/json-iterator/go"
	. "github.com/luvx21/coding-go/coding-common/common_x/alias_x"
	"github.com/luvx21/coding-go/coding-usage/api/common"
	"github.com/tidwall/gjson"
)

type User common.User

var users = [2]User{{1, "foo", 18}, {2, "bar", 19}}

func Test_00(t *testing.T) {
	// 序列化
	jsonBlob, _ := json.MarshalIndent(users, "", "    ")
	fmt.Println(string(jsonBlob))

	// 反序列化
	_ = json.Unmarshal(jsonBlob, &users)
	fmt.Printf("%+v\n", users)

	// 不反序列化为对象, 直接操作
	var f SliceAny
	//ff := make(JsonObject)
	_ = json.Unmarshal(jsonBlob, &f)

	for k, v := range f {
		if kvs, ok := v.(JsonObject); ok {
			fmt.Println(k, "Id: ", kvs["Id"], "name: ", kvs["name"])
		}
	}

	// var f any
	// jsonBlob := []byte(`{"Id":1 ,"Name":"foo", "Parents":["Gomez", "Morticia"]}`)
	// json.Unmarshal(jsonBlob, &f)
	// for k, v := range f.(JsonObject) {
	// }
}

func Test_b(t *testing.T) {
	body := `{
    "data": {
        "body": "aaaa"
    },
    "aa": 1
}`
	m := make(JsonObject)
	_ = json.Unmarshal([]byte(body), &m)

	a2 := m["data"]
	if kvs, ok := a2.(JsonObject); ok {
		fmt.Println(a2, kvs["body"])
	}
}

func Test_sonic_00(t *testing.T) {
	_json, _ := sonic.Marshal(&users)

	var data []common.User
	_ = sonic.Unmarshal(_json, &data)

	fmt.Println(string(_json), data)
}

func Test_JsonIter_00(t *testing.T) {
	var json = goJson.ConfigCompatibleWithStandardLibrary

	jsonBlob, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(jsonBlob))

	var data []common.User
	err = json.Unmarshal(jsonBlob, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", data)

}

func Test_gjson_00(t *testing.T) {
	_json := `
{
    "a": {
        "b": "bbb"
    }
}
`
	s := gjson.Get(_json, "a.b").String()
	fmt.Println(s)
}
