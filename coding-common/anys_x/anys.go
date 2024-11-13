package anys_x

import (
    "fmt"
)

func String[T any](s T) string {
    return fmt.Sprintf("%v", s)
}
