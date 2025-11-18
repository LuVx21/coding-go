package optional_x

import (
	"fmt"

	"github.com/luvx21/coding-go/coding-common/common_x/funcs"
	"github.com/luvx21/coding-go/coding-common/reflects"
)

type Optional[T any] struct {
	data T
}

func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

func Of[T any](value T) Optional[T] {
	return Optional[T]{data: value}
}

func OfNullable[T any](value T) Optional[T] {
	if reflects.IsNil(value) {
		return Empty[T]()
	}
	return Of(value)
}

func (op Optional[T]) Get() T {
	return op.data
}

func (op Optional[T]) IsPresent() bool {
	return !reflects.IsNil(op.data)
}

func (op Optional[T]) IsEmpty() bool {
	return !op.IsPresent()
}

func (op Optional[T]) IfPresent(consumer funcs.Consumer[T]) {
	if op.IsPresent() {
		consumer(op.data)
	}
}

func (op Optional[T]) IfPresentOrElse(consumer funcs.Consumer[T], run funcs.Runnable) {
	if op.IsPresent() {
		consumer(op.data)
	} else {
		run()
	}
}

func (op Optional[T]) Filter(predicate funcs.Predicate[T]) Optional[T] {
	if op.IsEmpty() {
		return op
	}
	if predicate(op.data) {
		return op
	}
	return Empty[T]()
}

func (op Optional[T]) Map(mapper funcs.Function[T, T]) Optional[T] {
	if op.IsEmpty() {
		return Empty[T]()
	}
	return OfNullable(mapper(op.data))
}

func (op Optional[T]) FlatMap(mapper funcs.Function[T, Optional[T]]) Optional[T] {
	if op.IsEmpty() {
		return Empty[T]()
	}
	return mapper(op.data)
}

func (op Optional[T]) Or(supplier funcs.Supplier[Optional[T]]) Optional[T] {
	if op.IsPresent() {
		return op
	}
	return supplier()
}

func (op Optional[T]) OrElse(value T) T {
	if op.IsPresent() {
		return op.data
	}
	return value
}

func (op Optional[T]) OrElseGet(supplier funcs.Supplier[T]) T {
	if op.IsPresent() {
		return op.data
	}
	return supplier()
}

func (op Optional[T]) String() string {
	if op.IsPresent() {
		return fmt.Sprintf("%v", op.data)
	} else {
		return "Optional.empty"
	}
}

//----------

//func (op Optional[T]) Equals(value T) bool {
//    if op.IsPresent() {
//        return op.data == value
//    }
//    return false
//}
