package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/wire"
)

var (
	ProviderSet = wire.NewSet(NewX, NewY, NewZ)
)

type (
	X struct{ Value int }
	Y struct{ Value int }
	Z struct{ Value int }
)

func NewX() X { return X{Value: 7} }
func NewY(x X) Y {
	fmt.Println("入参x", x.Value)
	return Y{Value: x.Value + 1}
}
func NewZ(ctx context.Context, y Y) (Z, error) {
	fmt.Println("入参y", y.Value)
	if y.Value == 0 {
		return Z{}, errors.New("cannot provide z when value is zero")
	}
	return Z{Value: y.Value + 2}, nil
}
