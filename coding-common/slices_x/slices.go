package slices_x

import (
    "github.com/luvx21/coding-go/coding-common/common_x"
    . "github.com/luvx21/coding-go/coding-common/common_x/types_x"
    "github.com/luvx21/coding-go/coding-common/reflects"
)

func First[S ~[]E, E any](s S) (E, bool) {
    if len(s) == 0 {
        var zero E
        return zero, false
    }
    return s[0], true
}

func Last[S ~[]E, E any](s S) (E, bool) {
    l := len(s)
    if l == 0 {
        var zero E
        return zero, false
    }
    return s[l-1], true
}

func Transfer[I, O any](f Function[I, O], s ...I) []O {
    r := make([]O, len(s))
    for i, v := range s {
        r[i] = f(v)
    }
    return r
}

// ToAnySliceE 入参类型一致
func ToAnySliceE[T any](s ...T) []any {
    f := func(a T) any { return a }
    return Transfer[T, any](f, s...)
}

// ToAnySlice 入参类型可随意
func ToAnySlice(s ...any) []any {
    return ToAnySliceE[any](s...)
}

func IsEmpty[S ~[]E, E any](s S) (bool, S) {
    if s == nil || len(s) == 0 {
        return true, s
    }
    return false, s
}

func ClearZero[S ~[]E, E comparable](s S) S {
    r := make([]E, 0, len(s))
    for i := range s {
        if !common_x.IsZero(s[i]) {
            r = append(r, s[i])
        }
    }
    return r
}

func ClearZeroRef[S ~[]E, E any](s S) S {
    r := make([]E, 0, len(s))
    for i := range s {
        if !reflects.IsZeroRef(s[i]) {
            r = append(r, s[i])
        }
    }
    return r
}
