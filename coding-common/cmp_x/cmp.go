package cmp_x

import "cmp"

func Between[T cmp.Ordered](value, min, max T) bool {
	return value > min && value < max
}

func BetweenAnd[T cmp.Ordered](v, min, max T) bool {
	return v >= min && v <= max
}

func CmpBy[T any](by func(T, T) bool) func(T, T) int {
	return func(a, b T) int {
		if by(a, b) {
			return -1
		} else if by(b, a) {
			return 1
		}
		return 0
	}
}
