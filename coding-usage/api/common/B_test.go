package common

import (
	"fmt"
	"testing"
)

func Test01(t *testing.T) {
	var b MyInt
	fmt.Println(b.IsZero())
	b = 1
	fmt.Println(b.Add(2))
}
