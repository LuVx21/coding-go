package tree

import (
	"fmt"
	"testing"
)

func Test_string(t *testing.T) {
	root := NewCBT(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 100)
	fmt.Println(root)
}
