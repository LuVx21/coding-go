package jsons

import (
    "github.com/bytedance/sonic"
)

func JsonStringToObject(s string, v any) error {
    return sonic.Unmarshal([]byte(s), v)
}

func JsonStringToMap[K comparable, V any, M ~map[K]V](s string) (M, error) {
    m := make(M)
    err := sonic.Unmarshal([]byte(s), &m)
    return m, err
}

func JsonStringToArray[E any, S ~[]E](s string) (S, error) {
    slice := make(S, 0)
    err := sonic.Unmarshal([]byte(s), &slice)
    return slice, err
}
