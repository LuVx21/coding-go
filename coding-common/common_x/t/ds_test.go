package t

import (
	"fmt"
	"slices"
	"testing"
)

func Test_ds_list_00(t *testing.T) {
	lst := List[int]{}
	lst.Push(10, 13, 23)

	for e := range lst.All() {
		fmt.Println("元素:", e)
	}

	all := slices.Collect(lst.All())
	fmt.Println("all:", all)
}
