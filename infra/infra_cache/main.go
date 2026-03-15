package infra_cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/luvx21/coding-go/coding-common/cast_x"
)

// BigCacheDataRetrieval K 能转string
type BigCacheDataRetrieval[K comparable, V []byte] struct {
	cache *bigcache.BigCache
}

// Get implements [cache.MultiDataRetrievable].
func (s *BigCacheDataRetrieval[K, V]) Get(ids []K) map[K]V {
	if s.cache == nil || len(ids) == 0 {
		return nil
	}
	r := make(map[K]V, len(ids))
	for _, id := range ids {
		_id, err := cast_x.ToStringE(id)
		if err != nil {
			slog.Warn("非法cache key,需String()")
			continue
		}
		if bs, err := s.cache.Get(_id); err == nil {
			r[id] = bs
		}
	}
	return r
}

// Set implements [cache.MultiDataRetrievable].
func (s *BigCacheDataRetrieval[K, V]) Set(data map[K]V) {
	if s.cache == nil {
		c, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
		if err != nil {
			slog.Error("初始化缓存错误", "err", err)
			return
		}
		s.cache = c
	}
	for k, v := range data {
		_id, err := cast_x.ToStringE(k)
		if err != nil {
			slog.Warn("非法cache key,需String()")
			continue
		}
		s.cache.Set(_id, v)
	}
}

type keyType interface {
	ristretto.Key
	comparable
}

type RistrettoDataRetrieval[K keyType, V []byte] struct {
	cache *ristretto.Cache[K, V]
}

// Get implements [cache.MultiDataRetrievable].
func (s *RistrettoDataRetrieval[K, V]) Get(ids []K) map[K]V {
	if s.cache == nil || len(ids) == 0 {
		return nil
	}
	r := make(map[K]V, len(ids))
	for _, id := range ids {
		if bs, exist := s.cache.Get(id); exist {
			r[id] = bs
		}
	}
	return r
}

// Set implements [cache.MultiDataRetrievable].
func (s *RistrettoDataRetrieval[K, V]) Set(data map[K]V) {
	if s.cache == nil {
		c, err := ristretto.NewCache(&ristretto.Config[K, V]{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		})
		if err != nil {
			slog.Error("初始化缓存错误", "err", err)
			return
		}
		s.cache = c
	}
	for k, v := range data {
		s.cache.Set(k, v, 1)
	}
}
