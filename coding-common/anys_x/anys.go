package anys_x

import (
	"fmt"
)

func String[T any](s T) string {
	return fmt.Sprintf("%v", s)
}

func NilGet[T any](t *T, f func() *T) *T {
	if t != nil {
		return t
	}
	return f()
}
