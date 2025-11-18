package fnv_x

import (
	"fmt"
	"testing"
)

func Test_fnv_00(t *testing.T) {
	i := Fnv32("hello")
	fmt.Println(i)
}
