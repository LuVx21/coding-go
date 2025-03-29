package msgpack

import (
	"fmt"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

type user struct {
	Name string
	Age  int32
}

func Test_msgpack_00(t *testing.T) {
	u := user{Name: "foobar", Age: 18}
	bytes, _ := msgpack.Marshal(u)

	var r user
	msgpack.Unmarshal(bytes, &r)
	fmt.Println(r.Name, r.Age)
}
