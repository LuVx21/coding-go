package maps_x

import "github.com/luvx21/coding-go/coding-common/reflects"

func RemoveIf[M ~map[K]V, K comparable, V any](m M, f func(K, V) bool) {
    for k, v := range m {
        if f(k, v) {
            delete(m, k)
        }
    }
}

func GetOrDefault[M ~map[K]V, K comparable, V any](m M, k K, defaultV V) V {
    if v, exist := m[k]; exist {
        return v
    } else {
        return defaultV
    }
}

func Compute[M ~map[K]V, K comparable, V any](m M, k K, f func(K, V) V) V {
    var zero V
    oldValue, exist := m[k]
    newValue := f(k, oldValue)

    if reflects.IsNil(newValue) {
        if !reflects.IsNil(oldValue) || exist {
            delete(m, k)
            return zero
        } else {
            return zero
        }
    } else {
        m[k] = newValue
        return newValue
    }
}

func ComputeIfAbsent[M ~map[K]V, K comparable, V any](m M, k K, f func(K) V) V {
    oldValue := m[k]
    if reflects.IsNil(oldValue) {
        newValue := f(k)
        if !reflects.IsNil(newValue) {
            m[k] = newValue
            return newValue
        }
    }
    return oldValue
}

func ComputeIfPresent[M ~map[K]V, K comparable, V any](m M, k K, f func(K, V) V) V {
    var zero V
    oldValue := m[k]
    if !reflects.IsNil(oldValue) {
        newValue := f(k, oldValue)
        if !reflects.IsNil(newValue) {
            m[k] = newValue
            return newValue
        } else {
            delete(m, k)
            return zero
        }
    } else {
        return oldValue
    }
}
