package weibo_p

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"luvx/gin/common/consts"

	"github.com/bytedance/sonic"
	"github.com/luvx21/coding-go/coding-common/common_x/types_x"
	"github.com/luvx21/coding-go/coding-common/nets_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/retry"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

func requestWeibo(url string, queryMap map[string]any, headerMap map[string]string) (gorequest.Response, string, []error) {
	pUrl, _ := nets_x.UrlAddQuery(url, queryMap)
	consts.GetRateLimiter(url).Wait(context.TODO())

	gg := consts.GoRequest.Get(pUrl.String())

	defaultHeader := map[string]string{
		"User-Agent": consts.UserAgent,
		"Host":       "weibo.com",
		"Referer":    "https://weibo.com/mygroups?gid=4670120389774996",
	}
	for k, v := range defaultHeader {
		gg.Set(k, v)
	}
	for k, v := range headerMap {
		gg.Set(k, v)
	}

	t, _ := retry.SupplyWithRetry("weibo请求重试", func() types_x.Tuple[gorequest.Response, string, []error] {
		r, body, errs := gg.End()

		if len(errs) > 0 || r.StatusCode/100 != 2 {
			log.Errorln("weibo请求异常", url, errs, r.Status)
			panic("fast-fail retry: weibo请求异常")
		}

		isJson := sonic.ValidString(body)
		slog.Debug("weibo请求信息", "请求", pUrl, "响应", r.StatusCode, "Json", isJson)
		if !isJson {
			msg := "weibo->请求结果非json,cookie可能过期"
			slog.Warn(msg, "响应空", r == nil, "异常", errs, "url", pUrl.String())
			os_x.Command("kv", "set", "weibo_cookie", "1")
			// redis_dao.SetSwitch("remote_cookie", true)
			// panic("fast-fail retry: " + msg)
			return types_x.NewTuple[gorequest.Response](nil, "", []error{fmt.Errorf("%s", msg)})
		}

		return types_x.NewTuple(r, body, errs)
	}, 5, 5*time.Second)

	return t.A, t.B, t.C
}
