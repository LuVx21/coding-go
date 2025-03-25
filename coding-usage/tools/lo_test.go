package tools

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func Test_a(t *testing.T) {
	names := lo.Uniq([]string{"foo", "bar", "foo"})
	fmt.Println(names)

	lo.ForEach([]string{"hello", "world"}, func(x string, _ int) {
		println(x)
	})

	groups := lo.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
		return i % 3
	})
	fmt.Println(groups)
}
