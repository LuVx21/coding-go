package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/golang/groupcache"
)

func Test_group_cache(t *testing.T) {
	// 1. 创建 Groupcache 实例
	pool := groupcache.NewHTTPPoolOpts("http://localhost:8080", &groupcache.HTTPPoolOptions{})

	// 2. 定义缓存组
	var cacheGroup *groupcache.Group
	cacheGroup = groupcache.NewGroup("example-group", 64<<20, groupcache.GetterFunc(
		func(ctx context.Context, key string, dest groupcache.Sink) error {
			// 当缓存未命中时，从数据源获取数据（此处模拟从数据库或其他来源）
			fmt.Printf("缓存未命中, 加载缓存, key: %s\n", key)
			dest.SetString("value-for-" + key) // 设置缓存值
			return nil
		},
	))

	// 3. 启动 HTTP 服务（用于节点间通信）
	go func() {
		http.Handle("/_groupcache/", pool)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// 4. 模拟从缓存中获取数据
	var result string
	ctx := context.Background()
	key := "foobar"
	if err := cacheGroup.Get(ctx, key, groupcache.StringSink(&result)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("缓存值:", result)

	// 再次获取相同 key，直接从缓存读取
	if err := cacheGroup.Get(ctx, key, groupcache.StringSink(&result)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("缓存值 (cached):", result)
}
