package anys_x

func AnyNil[T any](as ...T) bool {
    for _, a := range as {
        if a == nil {
            return true
        }
    }
    return false
}
