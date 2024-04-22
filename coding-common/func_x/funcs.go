package func_x

func Identity[T any]() func(T) T {
    return func(t T) T { return t }
}
