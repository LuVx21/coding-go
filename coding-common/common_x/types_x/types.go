package types_x

type Map[K comparable, V any] map[K]V
type Set[E comparable] Map[E, struct{}]

type Function[T, R any] func(T) R
type Consumer[T any] func(T)
type Supplier[T any] func() T
type Runnable func()
type Predicate[T any] Function[T, bool]

type BiFunction[I1, I2, R any] func(I1, I2) R
type BiConsumer[I1, I2 any] func(I1, I2)
type BiPredicate[I1, I2 any] BiFunction[I1, I2, bool]

func (m *Map[K, V]) Filter(f func(K, V) bool) Map[K, V] {
    clone := make(Map[K, V], len(*m))
    for k, v := range *m {
        if f(k, v) {
            clone[k] = v
        }
    }
    return clone
}

func (m *Map[K, V]) Clone() Map[K, V] {
    f := func(k K, v V) bool { return true }
    return m.Filter(f)
}

func (m *Map[K, V]) Merge(source Map[K, V], replace bool) {
    if *m == nil {
        *m = make(Map[K, V], len(source))
    }

    for sourceKey, sourceValue := range source {
        if _, ok := (*m)[sourceKey]; !ok || replace {
            (*m)[sourceKey] = sourceValue
        }
    }
}

func (s *Set[E]) Add(e ...E) {
    for _, _e := range e {
        (*s)[_e] = struct{}{}
    }
}

func (s *Set[E]) Remove(e ...E) {
    for _, _e := range e {
        delete(*s, _e)
    }
}

func (s *Set[E]) Contain(e E) bool {
    _, exist := (*s)[e]
    return exist
}
