package consts

import (
    "github.com/luvx21/coding-go/coding-common/ids"
    "github.com/parnurzeal/gorequest"
    "golang.org/x/sync/singleflight"
    "golang.org/x/time/rate"
)

const (
    UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
)

var (
    GoRequest   = gorequest.New()
    RateLimiter = rate.NewLimiter(1, 1)
    IdWorker, _ = ids.NewSnowflakeIdWorker(0, 0)
    SfGroup     = singleflight.Group{}
)
