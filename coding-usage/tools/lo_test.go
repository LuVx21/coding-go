package tools

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/samber/lo"
	lp "github.com/samber/lo/parallel"
)

func Test_a(t *testing.T) {
	names := lo.Uniq([]string{"foo", "bar", "foo"})
	fmt.Println(names)

	lo.ForEach([]string{"hello", "world"}, func(x string, _ int) {
		println(x)
	})
}

func Test_parallel_00(t *testing.T) {
	groups := lo.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int { return i % 3 })
	fmt.Println(groups)

	s1 := lp.PartitionBy([]int{1, 2, 3, 4, 2}, func(n int) string { return "_" + strconv.Itoa(n) })
	fmt.Println(s1)
}
