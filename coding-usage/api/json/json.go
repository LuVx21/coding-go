package json

import (
    "encoding/json"
    "fmt"
    goJson "github.com/json-iterator/go"
    "github.com/luvx21/coding-go/coding-usage/api/common"
    "github.com/pkg/errors"
)

type User common.User

var users = [2]User{{1, "foo", 18}, {2, "bar", 19}}

func a() []byte {
    // 序列化
    jsonBlob, err := json.MarshalIndent(users, "", "    ")
    if err != nil {
        err = errors.Wrap(err, "序列化失败")
        fmt.Println("error: ", err)
    }
    fmt.Println(string(jsonBlob))

    // 反序列化
    err = json.Unmarshal(jsonBlob, &users)
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Printf("%+v\n", users)

    return jsonBlob
}

func b() []byte {
    var json = goJson.ConfigCompatibleWithStandardLibrary

    jsonBlob, err := json.MarshalIndent(users, "", "    ")
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Println(string(jsonBlob))

    err = json.Unmarshal(jsonBlob, &users)
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Printf("%+v\n", users)

    return jsonBlob
}

func c() {
    jsonBlob := a()
    // 不反序列化为对象, 直接操作
    var f []interface{}
    //ff := make(map[string]interface{})
    json.Unmarshal(jsonBlob, &f)

    for k, v := range f {
        if kvs, ok := v.(map[string]interface{}); ok {
            fmt.Println(k, "Id: ", kvs["Id"], "name: ", kvs["name"])
        }
    }

    // var f interface{}
    // jsonBlob := []byte(`{"Id":1 ,"Name":"foo", "Parents":["Gomez", "Morticia"]}`)
    // json.Unmarshal(jsonBlob, &f)
    // for k, v := range f.(map[string]interface{}) {
    // }
}
