package optional_x

import (
	"fmt"
	"testing"
)

type test struct {
	a string
}

func Test_01(t *testing.T) {
	var a int
	present := OfNullable(a).
		Map(func(i int) int {
			return i + 100
		}).
		OrElseGet(func() int {
			return 999
		})
	fmt.Println(present, a, a == 0)

	var t2 test
	isPresent := OfNullable(t2).
		IsEmpty()
	fmt.Println(isPresent, len(t2.a))
}
