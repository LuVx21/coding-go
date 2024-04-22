package db

import (
    "fmt"
    "github.com/redis/go-redis/v9"
    "luvx/gin/config"
)

func NewRedisClient() *redis.Client {
    fmt.Println("初始化Redis连接...")
    redisConfig := config.AppConfig.Redis
    return redis.NewClient(&redis.Options{
        Addr:     redisConfig.Host,
        Username: redisConfig.Username,
        Password: redisConfig.Password,
        DB:       0,
    })
}
