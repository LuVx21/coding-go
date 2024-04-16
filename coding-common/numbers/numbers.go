package numbers

import "golang.org/x/exp/constraints"

func BetweenAnd[T constraints.Integer | constraints.Float](v, start, end T) bool {
    return v >= start && v <= end
}
