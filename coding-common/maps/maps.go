package maps

func GetOrDefault[K comparable, V any](m map[K]V, k K, defaultV V) V {
    if v, exist := m[k]; exist {
        return v
    } else {
        return defaultV
    }
}

func Compute[K comparable, V comparable](m map[K]V, k K, f func(K, V) V) V {
    var zero V
    oldValue, exist := m[k]
    newValue := f(k, oldValue)

    if newValue == zero {
        if oldValue != zero || exist {
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

func ComputeIfAbsent[K comparable, V comparable](m map[K]V, k K, f func(K) V) V {
    var zero V
    oldValue := m[k]
    if oldValue == zero {
        newValue := f(k)
        if newValue != zero {
            m[k] = newValue
            return newValue
        }
    }
    return oldValue
}

func ComputeIfPresent[K comparable, V comparable](m map[K]V, k K, f func(K, V) V) V {
    var zero V
    oldValue := m[k]
    if oldValue != zero {
        newValue := f(k, oldValue)
        if newValue != zero {
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
