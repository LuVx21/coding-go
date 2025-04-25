package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	go_cache "github.com/patrickmn/go-cache"
)

func Test_gocache(t *testing.T) {
	c := go_cache.New(5*time.Minute, 10*time.Minute)
	c.Set("foo", "bar", go_cache.DefaultExpiration)
	v, exist := c.Get("foo")
	fmt.Println(v, exist)
}

func Test_ristretto(t *testing.T) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	cache.SetWithTTL("key", "value", 1, time.Hour)
	cache.Wait()
	value, found := cache.Get("key")
	if !found {
		panic("missing value")
	}
	fmt.Println(value)

	// del value from cache
	cache.Del("key")
}
