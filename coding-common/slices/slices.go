package slices

import (
    "github.com/luvx21/coding-go/coding-common/common"
    "github.com/luvx21/coding-go/coding-common/reflects"
)

func ClearZero[S ~[]E, E comparable](s S) S {
    r := make([]E, 0, len(s))
    for i := range s {
        if !common.IsZero(s[i]) {
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
