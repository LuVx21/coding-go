package infra_cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/luvx21/coding-go/coding-common/cache"
)

func Test_cache_00(t *testing.T) {
	bc1, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	ds1 := &BigCacheDataRetrieval[string, []byte]{cache: bc1}
	bc2, _ := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	bc2.Set("1", []byte("aa"), 1)
	ds2 := &RistrettoDataRetrieval[string, []byte]{cache: bc2}

	m := cache.GetFrom([]string{"1", "2"}, ds1, ds2)
	for k, v := range m {
		fmt.Println(k, "=", string(v))
	}
}
