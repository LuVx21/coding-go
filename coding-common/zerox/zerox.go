package zerox

import "github.com/luvx21/coding-go/coding-common/reflects"

func ZeroGet[T any](t T, f func() T) T {
    if !reflects.IsZeroRef(t) {
        return t
    }
    return f()
}
