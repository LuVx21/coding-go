package maths_x

import (
	"cmp"

	"github.com/luvx21/coding-go/coding-common/common_x/types_x"
)

func Abs[T types_x.Number](val T) T {
	if val < 0 {
		return -val
	}
	return val
}

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
