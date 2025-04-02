package consts

import (
	"os"
	"time"

	"github.com/luvx21/coding-go/coding-common/ids"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"
)

const (
	ServiceHost = "mini.rx"
	ImgHost     = "img.rx"
	UserAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	AppSwitchKey = "app_switch"
)

var (
	GoRequest   = gorequest.New().Timeout(time.Minute)
	RateLimiter = rate.NewLimiter(1, 1)
	IdWorker, _ = ids.NewSnowflakeIdWorker(0, 0)
	SfGroup     = singleflight.Group{}
)

var (
	AppPort     = os.Getenv("APP_PORT")
	AppHostName = os.Getenv("APP_HOST_NAME")
	AppProxy    = os.Getenv("APP_PROXY")
	PasswordStr = os.Getenv("passwordStr")
)
