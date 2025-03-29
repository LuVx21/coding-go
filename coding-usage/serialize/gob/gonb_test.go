package gob

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

type student struct {
	Name string
	Age  int32
}

func Test_gob_00(t *testing.T) {
	studentEncode := student{Name: "foobar", Age: 30}
	var b bytes.Buffer
	_ = gob.NewEncoder(&b).Encode(studentEncode)

	var r student
	_ = gob.NewDecoder(&b).Decode(&r)

	fmt.Println(r.Name, r.Age)
}
