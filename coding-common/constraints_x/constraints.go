package constraints_x

type Map[K comparable, V any] interface{ ~map[K]V }
