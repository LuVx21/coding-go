package probability

import (
	"math/rand/v2"
)

// Percent 以 n% 的概率返回 true，并发安全
func Percent(n int) bool {
	if n <= 0 {
		return false
	}
	if n >= 100 {
		return true
	}
	return rand.IntN(100) < n
}
