//go:build wireinject
// +build wireinject

package main

import (
    "context"
    "fmt"
    "github.com/google/wire"
)

var ZZ Z

func initZ(ctx context.Context) (Z, error) {
    wire.Build(ProviderSet)
    return Z{}, nil
}

func main() {
    ZZ, _ := initZ(context.TODO())
    fmt.Println(ZZ.Value)
}
