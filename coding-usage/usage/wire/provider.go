package main

import (
	"context"
	"errors"
	"github.com/google/wire"
)

type X struct {
	Value int
}

// NewX 返回一个X对象
func NewX() X {
	return X{Value: 7}
}

type Y struct {
	Value int
}

// NewY 返回一个Y对象，需要传入一个X对象作为依赖。
func NewY(x X) Y {
	return Y{Value: x.Value + 1}
}

type Z struct {
	Value int
}

// NewZ 返回一个Z对象，当传入依赖的value为0时会返回错误。
func NewZ(ctx context.Context, y Y) (Z, error) {
	if y.Value == 0 {
		return Z{}, errors.New("cannot provide z when value is zero")
	}
	return Z{Value: y.Value + 2}, nil
}

var ProviderSet = wire.NewSet(NewX, NewY, NewZ)
