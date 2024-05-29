package maths_x

import "golang.org/x/exp/constraints"

func Abs[T constraints.Integer | constraints.Float](val T) T {
    if val < 0 {
        return -val
    }
    return val
}
