package bili

import (
	"context"
	"luvx/gin/common/consts"
	"luvx/gin/service/cookie"

	"github.com/bytedance/sonic"
	"github.com/luvx21/coding-go/coding-common/nets_x"
	log "github.com/sirupsen/logrus"
)

var ()

func biliRequest(_url string, queryMap map[string]any, useCookie bool) string {
	pUrl, _ := nets_x.UrlAddQuery(_url, queryMap)

	consts.GetRateLimiter(_url).Wait(context.TODO())
	log.Infoln("请求:", pUrl)
	sa := consts.GoRequest.Get(pUrl.String()).
		Set("User-Agent", consts.UserAgent).
		Set("Referer", "https://www.bilibili.com/")
	if useCookie {
		sa = sa.Set("Cookie", cookie.GetCookieStrByHost(".bilibili.com"))
	}
	r, body, errs := sa.End()

	if !sonic.ValidString(body) {
		log.Warnln("bili->请求结果非json,cookie可能过期", r == nil, body, errs)
		return ""
	}

	return body
}
