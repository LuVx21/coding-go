package slices_x

import (
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/reflects"
)

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
