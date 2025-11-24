package maps_x

func MapAny[M ~map[K]V, K comparable, V any](m M) map[K]any {
	r := make(map[K]any, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}

func TypedMap[M ~map[K]any, K comparable, V any](m M) map[K]V {
	r := make(map[K]V, len(m))
	for k, v := range m {
		if val, ok := v.(V); ok {
			r[k] = val
		}
	}
	return r
}
