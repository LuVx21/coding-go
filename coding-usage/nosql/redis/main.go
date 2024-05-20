package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/luvx21/coding-go/coding-usage/nosql"
    "github.com/redis/go-redis/v9"
)

var (
    addr     = nosql.RedisAddr
    userName = nosql.RedisUserName
    password = nosql.RedisPassword
    db       = 0
)

func connect() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     addr,
        Username: userName,
        Password: password,
        DB:       db,
    })
}

func ExampleClient() {
    rdb := connect()

    var ctx = context.Background()
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
