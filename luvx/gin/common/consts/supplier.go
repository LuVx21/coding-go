package consts

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	gocache "github.com/eko/gocache/lib/v4/cache"
	bigcachestore "github.com/eko/gocache/store/bigcache/v4"
	"github.com/google/uuid"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"golang.org/x/time/rate"
)

var (
	rateLimiterMu  sync.Mutex
	rateLimiterMap = maps_x.Map[string, *rate.Limiter]{}
	onceMu         sync.Mutex
	onceMap        = maps_x.Map[string, *sync.Once]{}
)

func NewLoadableCache[T any](loadFunc gocache.LoadFunction[T]) *gocache.LoadableCache[T] {
	bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	bigcacheStore := bigcachestore.NewBigcache(bigcacheClient)
	return gocache.NewLoadable(loadFunc, gocache.New[T](bigcacheStore))
}

func UUID() string {
	value := uuid.New()
	return strings.ToLower(strings.ReplaceAll(value.String(), "-", ""))
}

func GetRateLimiter(_url string) *rate.Limiter {
	parse, _ := url.Parse(_url)

	rateLimiterMu.Lock()
	defer rateLimiterMu.Unlock()

	limiter := rateLimiterMap[parse.Host]
	if limiter == nil {
		limiter = rate.NewLimiter(1, 1)
		rateLimiterMap[parse.Host] = limiter
	}
	return limiter
}

func GetOnce(k string) *sync.Once {
	onceMu.Lock()
	defer onceMu.Unlock()

	once := onceMap[k]
	if once == nil {
		once = &sync.Once{}
		onceMap[k] = once
	}
	return once
}
