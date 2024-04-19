package consts

import (
    "context"
    "github.com/allegro/bigcache/v3"
    gocache "github.com/eko/gocache/lib/v4/cache"
    bigcachestore "github.com/eko/gocache/store/bigcache/v4"
    "time"
)

func NewLoadableCache[T any](loadFunc gocache.LoadFunction[T]) *gocache.LoadableCache[T] {
    bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
    bigcacheStore := bigcachestore.NewBigcache(bigcacheClient)
    return gocache.NewLoadable[T](loadFunc, gocache.New[T](bigcacheStore))
}
