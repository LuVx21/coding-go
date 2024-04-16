//go:build wireinject
// +build wireinject

package db

import (
    "github.com/google/wire"
    "github.com/redis/go-redis/v9"
)

func InitRedisClient() *redis.Client {
    wire.Build(NewRedisClient)
    return nil
}


