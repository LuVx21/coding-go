package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/allegro/bigcache/v3"
    "github.com/dgraph-io/ristretto"
    gocache "github.com/eko/gocache/lib/v4/cache"
    gocache_store "github.com/eko/gocache/lib/v4/store"
    gocache_store_bigcache "github.com/eko/gocache/store/bigcache/v4"
    gocache_store_redis "github.com/eko/gocache/store/redis/v4"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "github.com/luvx21/coding-go/coding-common/times_x"
    go_cache "github.com/patrickmn/go-cache"
    "github.com/redis/go-redis/v9"
    "log"
    "testing"
    "time"
)

func Test_gocache(t *testing.T) {
    c := go_cache.New(5*time.Minute, 10*time.Minute)
    c.Set("foo", "bar", go_cache.DefaultExpiration)
    v, exist := c.Get("foo")
    fmt.Println(v, exist)
}

func Test_ristretto(t *testing.T) {
    cache, err := ristretto.NewCache(&ristretto.Config{
        NumCounters: 1e7,
        MaxCost:     1 << 30,
        BufferItems: 64,
    })
    if err != nil {
        panic(err)
    }
    cache.Set("key", "value", 1)
    cache.Wait()
    value, found := cache.Get("key")
    if !found {
        panic("missing value")
    }
    fmt.Println(value)

    // del value from cache
    cache.Del("key")
}

func Test_bigcache(t *testing.T) {
    config := bigcache.Config{
        Shards:             1024,
        LifeWindow:         10 * time.Minute,
        CleanWindow:        5 * time.Minute,
        MaxEntriesInWindow: 1000 * 10 * 60,
        MaxEntrySize:       500,
        Verbose:            true,
        HardMaxCacheSize:   8192,
        OnRemove: func(k string, v []byte) {
            fmt.Println("已经删除", k, v)
        },
        OnRemoveWithReason: nil,
    }

    cache, initErr := bigcache.New(context.Background(), config)
    if initErr != nil {
        log.Fatal(initErr)
    }

    cache.Set("foo", []byte("bar"))

    for _, key := range []string{"foo", "foo1"} {
        if data, err := cache.Get(key); err == nil {
            fmt.Printf("key: %s 结果:%s\n", key, data)
        } else {
            fmt.Printf("key: %s 结果:%v\n", key, nil)
        }
    }
}

func Test_Gocache(t *testing.T) {
    bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
    cache := gocache.New[[]byte](gocache_store_bigcache.NewBigcache(bigcacheClient))

    loadFunc := func(_ context.Context, key any) ([]byte, []gocache_store.Option, error) {
        fmt.Println("自动加载缓存...", key)
        return []byte(cast_x.ToString(key) + "-" + times_x.TimeNow()), nil, nil
    }
    loadable := gocache.NewLoadable[[]byte](loadFunc, cache)

    if err := loadable.Set(context.TODO(), "foo", []byte("bar")); err != nil {
        panic(err)
    }

    for _, key := range []string{"foo", "foo1"} {
        if data, err := loadable.Get(context.TODO(), key); err == nil {
            fmt.Printf("命中缓存-> key: %s 结果:%s\n", key, data)
        }
    }

    time.Sleep(1 * time.Second)
    get, _ := loadable.Get(context.TODO(), "foo1")
    fmt.Println(string(get))
}

func Test_Gocache_redis(t *testing.T) {
    ctx := context.TODO()

    bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
    bigcacheStore := gocache_store_bigcache.NewBigcache(bigcacheClient)

    redisStore := gocache_store_redis.NewRedis(redis.NewClient(&redis.Options{
        Addr: "127.0.0.1:6379",
    }))

	cache := gocache.NewChain(
        gocache.New[any](bigcacheStore),
        gocache.New[any](redisStore),
    )

    key := "foo1"
    //if err := cache.Set(ctx, key, "my-value", gocache_store.WithExpiration(15*time.Second)); err != nil {
    //    panic(err)
    //}

    value, err := cache.Get(ctx, key)
    switch {
    case err == nil:
        fmt.Printf("命中缓存-> key: %s 结果:%s\n", key, value)
    case errors.Is(err, redis.Nil):
        fmt.Printf("未命中缓存-> key: %s\n", key)
    default:
        fmt.Printf("查询缓存失败-> %s: %v\n", key, err)
    }
}
