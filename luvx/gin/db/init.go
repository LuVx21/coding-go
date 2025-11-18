package db

import "github.com/redis/go-redis/v9"

var RedisClient *redis.Client

func init() {
	RedisClient = InitRedisClient()
}
