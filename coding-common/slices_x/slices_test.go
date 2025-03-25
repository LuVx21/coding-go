package slices_x

import (
	"fmt"
	"testing"
)

type user struct {
	names []string
}

func Test_Partition(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := Partition(ints, 5)
	fmt.Println(r)
}

func Test_01(t *testing.T) {
	ms := []map[string]string{
		{"foo": "foo"},
		nil,
		{"bar": "bar"},
	}
	fmt.Println(ClearZeroRef(ms))
}
func Test_02(t *testing.T) {
	r := FilterTransfer(func(i string) bool {
		return i == "a" || i == "c"
	}, func(i string) string {
		return i + "_1"
	}, "a", "b", "c")
	fmt.Println(r, len(r))
}

func Test_flatmap(t *testing.T) {
	users := []user{
		{names: []string{"a", "b"}},
		{names: []string{"c", "d"}},
	}

	names := FlatMap(users, func(u user) []string {
		return Transfer(func(s string) string { return s + s }, u.names...)
	})
	fmt.Println(names)
}

func Test_groupby(t *testing.T) {
	ints := []int{1, 2, 3, 1, 5, 1}
	m := GroupBy(ints,
		func(i int) int { return 2 * i },
		func(i int) int { return 3 * i },
	)
	fmt.Println(m)
}

func Test_empty_01(t *testing.T) {
	e := Empty[string]()
	fmt.Println(e)
}
