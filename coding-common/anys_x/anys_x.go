package anys_x

import (
    "fmt"
)

func AnyNil[T any](as ...T) bool {
    for _, a := range as {
        if a == nil {
            return true
        }
    }
    return false
}

func String[T any](s T) string {
    return fmt.Sprintf("%v", s)
}
