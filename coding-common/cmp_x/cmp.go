package cmp_x

import "cmp"

func Between[T cmp.Ordered](value, min, max T) bool {
	return value > min && value < max
}

func BetweenAnd[T cmp.Ordered](v, min, max T) bool {
	return v >= min && v <= max
}
