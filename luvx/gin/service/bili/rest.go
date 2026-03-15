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
		sa = sa.Set("Cookie", cookie.GetCookieStrByHost(".bilibili.com")[".bilibili.com"])
	}
	r, body, errs := sa.End()

	if len(errs) > 0 || r.StatusCode/100 != 2 || !sonic.ValidString(body) {
		log.Errorln("bili请求异常", _url, errs, r.Status, body)
		return ""
	}

	// ff := make(map[string]any)
	// _ = sonic.UnmarshalString(body, &ff)
	// if cast_x.ToInt32(ff["code"]) != 0 || ff["data"] == nil {
	// 	log.Errorln("bili->请求结果码非0", _url, errs, r.Status, body)
	// 	return ""
	// }

	return body
}
