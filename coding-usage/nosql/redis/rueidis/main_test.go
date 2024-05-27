package rueidis

import (
    "context"
    "fmt"
    "github.com/luvx21/coding-go/coding-usage/nosql"
    "github.com/redis/rueidis"
    "log"
    "testing"
)

var rdb rueidis.Client

func connect() rueidis.Client {
    client, err := rueidis.NewClient(rueidis.ClientOption{
        InitAddress:  []string{nosql.RedisAddr},
        Username:     nosql.RedisUserName,
        Password:     nosql.RedisPassword,
        DisableCache: true,
    })
    if err != nil {
        log.Fatal("连接异常:", err)
    }
    return client
}

func beforeAfter(caseName string) func() {
    if rdb == nil {
        rdb = connect()
    }
    return func() {
        fmt.Println(caseName, "end...")
    }
}

func Test_rueidis_00(t *testing.T) {
    defer beforeAfter("Test_rueidis_00")()

    build := rdb.B().Get().Key("foo").Build()
    val, _ := rdb.Do(context.TODO(), build).ToString()
    fmt.Println("foo", "=", val)
}

func Test_rueidis_01(t *testing.T) {
    defer beforeAfter("Test_rueidis_01")()

    ctx := context.Background()
    completed := rdb.B().Set().Key("key").Value("val").Nx().Build()
    err := rdb.Do(ctx, completed).Error()

    build := rdb.B().Hgetall().Key("hm").Build()
    hm, err := rdb.Do(ctx, build).AsStrMap()
    fmt.Println(hm, err)
}
