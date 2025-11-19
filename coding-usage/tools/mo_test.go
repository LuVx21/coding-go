package tools

import (
	"fmt"
	"testing"

	"github.com/samber/mo"
	"github.com/samber/mo/option"
)

func Test_mo_00(t *testing.T) {
	out := option.Pipe3(
		mo.Some(21),
		option.Map(func(v int) int { return v * 2 }),
		option.FlatMap(func(v int) mo.Option[int] { return mo.None[int]() }),
		option.Map(func(v int) int { return v + 21 }),
	)
	fmt.Println(out)
}
