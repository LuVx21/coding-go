package func_x

func Identity[T any](t T) T {
	return t
}

func True[T any](t T) bool {
	return true
}
func False[T any](t T) bool {
	return false
}
