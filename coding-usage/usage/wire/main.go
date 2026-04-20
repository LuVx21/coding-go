//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
)

func main() {
	var ZZ Z
	ZZ, _ = initZ(context.TODO())
	fmt.Println(ZZ.Value)
}
