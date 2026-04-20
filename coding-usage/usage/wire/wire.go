//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
)

func initZ(ctx context.Context) (Z, error) {
	wire.Build(ProviderSet)
	return Z{}, nil
}
