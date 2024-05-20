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

// Intersect 交集
func Intersect[S ~[]E, E comparable](a, b S) S {
    var r S
    mp := make(map[E]struct{}, len(a))
    for _, val := range a {
        mp[val] = struct{}{}
    }
    for _, val := range b {
        if _, ok := mp[val]; ok {
            r = append(r, val)
        }
    }
    return r
}

// Diff 差集a-b
func Diff[S ~[]E, E comparable](a, b S) S {
    var r S
    mp := make(map[E]struct{}, len(b))
    for _, val := range b {
        mp[val] = struct{}{}
    }
    for _, val := range a {
        if _, ok := mp[val]; !ok {
            r = append(r, val)
        }
    }
    return r
}

//Unique 切片去重实现
func Unique[S ~[]E, E comparable](arr S) S {
    r := make(S, 0, len(arr))
    mp := map[E]struct{}{}
    for _, e := range arr {
        if _, ok := mp[e]; !ok {
            mp[e] = struct{}{}
            r = append(r, e)
        }
    }
    return r
}
