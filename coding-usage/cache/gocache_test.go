package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gs_bigcache "github.com/eko/gocache/store/bigcache/v4"
	gs_redis "github.com/eko/gocache/store/redis/v4"
	gs_ristretto "github.com/eko/gocache/store/ristretto/v4"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/redis/go-redis/v9"
)

var (
	key, value   = "foo", "bar"
	keys, values = []byte(key), []byte(value)
)

func Test_Gocache_ristretto(t *testing.T) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		fmt.Println("xxx", err.Error())
	}
	ristrettoStore := gs_ristretto.NewRistretto(ristrettoCache)

	gocache := cache.New[string](ristrettoStore)
	err = gocache.Set(t.Context(), "foo", value)
	if err != nil {
		fmt.Println("set错误", err.Error())
	}
	// Ristretto异步写入
	// b := ristrettoCache.Set("foo", value, 2)
	// ristrettoCache.Wait()
	time.Sleep(1 * time.Second)

	d, _ := gocache.Get(t.Context(), key)
	fmt.Printf("命中缓存1-> key: %s 结果: %s\n", key, d)
}

func Test_Gocache_bigcache(t *testing.T) {
	bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(50*time.Minute))
	gocache := cache.New[string](gs_bigcache.NewBigcache(bigcacheClient))

	if err := gocache.Set(context.TODO(), key, value); err != nil {
		panic(err)
	}

	bytes, _ := bigcacheClient.Get(key)
	fmt.Printf("命中缓存-> key: %s 结果: %s\n", key, bytes)

	// 读不出来, 底层读取到的是字节数组, 实际泛型是string
	d, _ := gocache.GetCodec().Get(context.TODO(), key)
	fmt.Printf("命中缓存-> key: %s 结果: %s\n", key, d)

	// 读不出来数据
	d, _ = gocache.Get(t.Context(), key)
	fmt.Printf("命中缓存-> key: %s 结果: %s\n", key, d)
}

// 可自动加载
func Test_Gocache_Loadable(t *testing.T) {
	bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	gocache := cache.New[[]byte](gs_bigcache.NewBigcache(bigcacheClient))

	loadFunc := func(_ context.Context, key any) ([]byte, []store.Option, error) {
		fmt.Println("自动加载缓存...", key)
		return []byte(cast_x.ToString(key) + "-" + times_x.TimeNow()), nil, nil
	}
	loadable := cache.NewLoadable(loadFunc, gocache)

	if err := loadable.Set(context.TODO(), key, values); err != nil {
		panic(err)
	}

	for _, key := range []string{key, "foo1"} {
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
	bigcacheStore := gs_bigcache.NewBigcache(bigcacheClient)

	redisStore := gs_redis.NewRedis(redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	}))

	gocache := cache.NewChain(
		cache.New[any](bigcacheStore),
		cache.New[any](redisStore),
	)

	key := "foo1"
	// if err := cache.Set(ctx, key, "my-value", gocache_store.WithExpiration(15*time.Second)); err != nil {
	//    panic(err)
	// }

	value, err := gocache.Get(ctx, key)
	switch {
	case err == nil:
		fmt.Printf("命中缓存-> key: %s 结果:%s\n", key, value)
	case errors.Is(err, redis.Nil):
		fmt.Printf("未命中缓存-> key: %s\n", key)
	default:
		fmt.Printf("查询缓存失败-> %s: %v\n", key, err)
	}
}
