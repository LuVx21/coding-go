package a

import (
	"fmt"
	"testing"
)

func Test_alias_slice_00(t *testing.T) {
	a := S[int]{1, 2, 3}
	b := S[any]{1, 2, 3, "a"}

	c := AS{1, 2, 3, "a"}
	fmt.Println(a, b, c)

}

func Test_alias_map_01(t *testing.T) {
	m := M[uint, string]{1: "a", 2: "b"}
	fmt.Println(m)

	m1 := CAM[uint]{1: "a", 2: "b"}
	fmt.Println(m1)
}

func Test_table_01(t *testing.T) {
	r := make(Table[string], 20)
	kv, ok := r["a"]
	if !ok {
		kv = make(Row)
		r["a"] = kv
	}
	kv["aa"] = "aaa"
	fmt.Println(r)
}
