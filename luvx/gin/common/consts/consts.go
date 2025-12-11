package consts

import (
	"os"
	"time"

	"github.com/luvx21/coding-go/coding-common/ids"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0"

	AppSwitchKey = "app:switch"
)

const (
	Nbsp, Ensp = "&nbsp;", "&ensp;"
)

var (
	GoRequest   = gorequest.New().Timeout(time.Minute).Set("User-Agent", UserAgent)
	RateLimiter = rate.NewLimiter(1, 1)
	IdWorker, _ = ids.NewSnowflakeIdWorker(0, 0)
	SfGroup     = singleflight.Group{}
)

var (
	ImgHost, _ = os_x.LookupEnv("ImgHost", "img.rx")

	AppPort     = os.Getenv("APP_PORT")
	AppHostName = os.Getenv("APP_HOST_NAME")
	AppProxy    = os.Getenv("APP_PROXY")
	PasswordStr = os.Getenv("passwordStr")
)
