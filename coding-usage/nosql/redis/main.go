package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ExampleClient() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "redisConfig.Host",
        Username: "redisConfig.Username",
        Password: "redisConfig.Password",
        DB:       0,
    })

    key1, key2 := "key1", "key2"
    err := rdb.Set(ctx, key1, "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, key1).Result()
    if err != nil {
        panic(err)
    }
    fmt.Println(key1, "=", val)

    val2, err := rdb.Get(ctx, key2).Result()
    if errors.Is(err, redis.Nil) {
        fmt.Println("key2 不存在")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println(key2, "=", val2)
    }
}

func main() {
    ExampleClient()
}
