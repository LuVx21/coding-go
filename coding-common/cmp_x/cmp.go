package cmp_x

import (
	"cmp"
	"strconv"
	"strings"
)

func Between[T cmp.Ordered](value, min, max T) bool { return value > min && value < max }
func BetweenAnd[T cmp.Ordered](v, min, max T) bool  { return v >= min && v <= max }

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

// CompareVersion `major.minor.patch.pre`
// -1: old > new, 0:  new == old 1: old > new
// versionOrder: 字符串对数字, 版本高的对应的数字高(map[string]int{"preview": 0, "alpha": 1, "beta": 2, "rc": 3, "": 4})
func CompareVersion(old, new_ string, versionOrder map[string]int) int {
	old, new_ = strings.TrimPrefix(old, "v"), strings.TrimPrefix(new_, "v")
	oparts, nparts := strings.Split(old, "."), strings.Split(new_, ".")
	maxLen := max(len(oparts), len(nparts))
	for i := range maxLen {
		ov, nv := partValue(oparts, i, versionOrder), partValue(nparts, i, versionOrder)
		if ov == nv {
			continue
		}
		if nv > ov {
			return -1
		} else {
			return 1
		}
	}
	return 0
}
func partValue(parts []string, idx int, versionOrder map[string]int) int {
	if idx < 0 || idx >= len(parts) {
		return 0
	}
	s := parts[idx]
	if num, err := strconv.Atoi(s); err == nil {
		return num
	}
	if v, ok := versionOrder[strings.ToLower(s)]; ok {
		return v
	}
	return 100 + int(s[0])
}
