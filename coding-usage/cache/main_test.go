package main

import (
    "context"
    "fmt"
    "github.com/allegro/bigcache/v3"
    "github.com/dgraph-io/ristretto"
    gocache "github.com/eko/gocache/lib/v4/cache"
    bigcachestore "github.com/eko/gocache/store/bigcache/v4"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "github.com/luvx21/coding-go/coding-common/times_x"
    go_cache "github.com/patrickmn/go-cache"
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
    bigcacheStore := bigcachestore.NewBigcache(bigcacheClient)

    loadFunc := func(_ context.Context, key any) ([]byte, error) {
        fmt.Println("自动加载...", key)
        return []byte(cast_x.ToString(key) + "-" + times_x.TimeNow()), nil
    }
    loadable := gocache.NewLoadable[[]byte](loadFunc, gocache.New[[]byte](bigcacheStore))

    err := loadable.Set(context.TODO(), "foo", []byte("bar"))
    if err != nil {
        panic(err)
    }

    for _, key := range []string{"foo", "foo1"} {
        if data, err := loadable.Get(context.TODO(), key); err == nil {
            fmt.Printf("key: %s 结果:%s\n", key, data)
        }
    }

    time.Sleep(1 * time.Second)
    get, _ := loadable.Get(context.TODO(), "foo1")
    fmt.Println(string(get))
}
