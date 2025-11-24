package slices_x

var EmptySlice = []any{}

func Empty[T any]() []T {
	return []T{}
}
