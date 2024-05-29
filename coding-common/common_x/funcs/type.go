package funcs

type Function[T, R any] func(T) R
type Consumer[T any] func(T)
type Supplier[T any] func() T
type Runnable func()
type Predicate[T any] Function[T, bool]

type BiFunction[I1, I2, R any] func(I1, I2) R
type BiConsumer[I1, I2 any] func(I1, I2)
type BiPredicate[I1, I2 any] BiFunction[I1, I2, bool]
