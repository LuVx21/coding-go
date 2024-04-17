package constraints

type Map[K comparable, V any] interface{ ~map[K]V }
