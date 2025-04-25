package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
)

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
	cache.Delete("foo")
}
