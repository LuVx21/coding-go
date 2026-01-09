package consts

import (
	"os"
	"time"

	"github.com/luvx21/coding-go/coding-common/ids"
	"github.com/luvx21/coding-go/coding-common/strings_x"
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
	ImgRedirectUrlPrefix = strings_x.FirstNonEmpty(os.Getenv("IMG_REDIRECT_URL_PREFIX"), "https://image.baidu.com/search/down")

	AppPort       = os.Getenv("APP_PORT")
	ServiceDomain = os.Getenv("SERVICE_DOMAIN")
	AppProxy      = os.Getenv("APP_PROXY")
	PasswordStr   = os.Getenv("passwordStr")
)
