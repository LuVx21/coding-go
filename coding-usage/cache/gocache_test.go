package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dgraph-io/ristretto"
	gocache "github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/lib/v4/store"
	gocache_store_bigcache "github.com/eko/gocache/store/bigcache/v4"
	gocache_store_redis "github.com/eko/gocache/store/redis/v4"
	gocache_store_ristretto "github.com/eko/gocache/store/ristretto/v4"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/redis/go-redis/v9"
)

func a() *gocache_store_ristretto.RistrettoStore {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	return gocache_store_ristretto.NewRistretto(ristrettoCache)
}

func Test_Gocache_00(t *testing.T) {
	bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(50*time.Minute))
	cache := gocache.New[string](gocache_store_bigcache.NewBigcache(bigcacheClient))

	if err := cache.Set(context.TODO(), "foo", ("bar")); err != nil {
		panic(err)
	}

	bytes, err := bigcacheClient.Get("foo")
	fmt.Println("底层:", string(bytes), err)

	// 读不出来, 底层读取到的是字节数组, 实际泛型是string
	// cache.GetCodec().Get()
	d, _ := cache.GetCodec().Get(context.TODO(), "foo")
	fmt.Printf("命中缓存-> key: %s 结果: %s\n", "foo", d)
}

func Test_Gocache(t *testing.T) {
	bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	cache := gocache.New[[]byte](gocache_store_bigcache.NewBigcache(bigcacheClient))

	loadFunc := func(_ context.Context, key any) ([]byte, []gocache_store.Option, error) {
		fmt.Println("自动加载缓存...", key)
		return []byte(cast_x.ToString(key) + "-" + times_x.TimeNow()), nil, nil
	}
	loadable := gocache.NewLoadable(loadFunc, cache)

	if err := loadable.Set(context.TODO(), "foo", []byte("bar")); err != nil {
		panic(err)
	}

	for _, key := range []string{"foo", "foo1"} {
		if data, err := loadable.Get(context.TODO(), key); err == nil {
			fmt.Printf("命中缓存-> key: %s 结果: %s\n", key, data)
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
	// if err := cache.Set(ctx, key, "my-value", gocache_store.WithExpiration(15*time.Second)); err != nil {
	//    panic(err)
	// }

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
