package sortx

import (
    . "github.com/luvx21/coding-go/coding-common/common_x/types_x"
    "sort"
)

type SortWrapper[T any] struct {
    items []T
    By    BiPredicate[*T, *T]
}

func (sw *SortWrapper[T]) Len() int {
    return len(sw.items)
}
func (sw *SortWrapper[T]) Swap(i, j int) {
    sw.items[i], sw.items[j] = sw.items[j], sw.items[i]
}
func (sw *SortWrapper[T]) Less(i, j int) bool {
    return sw.By(&sw.items[i], &sw.items[j])
}

func Sort[T any](s []T, by BiPredicate[*T, *T]) {
    sort.Sort(&SortWrapper[T]{s, by})
}
