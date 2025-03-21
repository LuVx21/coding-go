package maps_x

type Map[K comparable, V any] map[K]V

func (m *Map[K, V]) Filter(f func(K, V) bool) Map[K, V] {
	return Filter(*m, f)
}

func (m *Map[K, V]) Clone() Map[K, V] {
	return Clone(*m)
}

func (m *Map[K, V]) Merge(source Map[K, V], replace bool) Map[K, V] {
	if *m == nil {
		*m = make(Map[K, V], len(source))
	}
	return Merge(*m, source, replace)
}
