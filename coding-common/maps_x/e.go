package maps_x

var EmptyMap = map[any]any{}

func Empty[K comparable, V any]() map[K]V {
	return map[K]V{}
}
