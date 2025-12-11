package webhook

import (
	"context"
	"fmt"
	"log/slog"
	"luvx/gin/common/consts"
	"luvx/gin/config"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/db"
	"strings"
	"time"

	"github.com/icloudza/fxjson"
	"github.com/luvx21/coding-go/coding-common/nets_x"
)

// https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index
const (
	url_token              = "https://api.weixin.qq.com/cgi-bin/stable_token"
	url_send_msg           = "https://api.weixin.qq.com/cgi-bin/message/template/send"
	skin_base_url          = "wxpushskin.yeyu0926.workers.dev"
	redis_key_access_token = "app:webhook:weixin:access_token"
)

var (
	TO_USER      = mongo_dao.DynamicConfig.Get()["wx_to_user"].(string)
	TEMPLATE_ID  = mongo_dao.DynamicConfig.Get()["template_id"].(string)
	tokenPayload = map[string]any{
		"grant_type":    "client_credential",
		"appid":         config.Viper.GetString("webhook.weixin.appId"),
		"secret":        config.Viper.GetString("webhook.weixin.secret"),
		"force_refresh": false,
	}
)

func getAccessToken() string {
	if token, err := db.RedisClient.Get(context.TODO(), redis_key_access_token).Result(); err == nil {
		return token
	}

	b := requestWeixin(url_token, tokenPayload)
	fj := fxjson.FromString(b)
	token, exp := fj.Get("access_token").StringOr(""), time.Second*time.Duration(fj.Get("expires_in").IntOr(0))
	if exp > 0 {
		db.RedisClient.Set(context.TODO(), redis_key_access_token, token, exp)
	}
	return token
}

type WeixinMsg struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

func SendMessage(touser, template_id, title string, msg map[string]WeixinMsg) (string, int64) {
	var sb strings.Builder
	for _, v := range msg {
		sb.WriteString(v.Value)
		sb.WriteString("\n")
	}
	sb.WriteString("-----------")

	pUrl, _ := nets_x.UrlAddQuery(url_send_msg, map[string]any{"access_token": getAccessToken()})
	payload := map[string]any{
		"touser":      touser,
		"template_id": template_id,
		"topcolor":    "#FF0000",
		"url":         fmt.Sprintf(`%s?title=%s&message=%s&date=%s`, skin_base_url, nets_x.EncodeURIComponent(title), nets_x.EncodeURIComponent(sb.String()), nets_x.EncodeURIComponent(time.Now().Format(time.DateTime))),
		"data":        msg,
	}

	_json := requestWeixin(pUrl.String(), payload)
	fx := fxjson.FromString(_json)
	return fx.Get("errmsg").StringOr("error"), fx.Get("msgid").IntOr(-1)
}

func requestWeixin(_url string, payload map[string]any) string {
	r, b, es := consts.GoRequest.
		Post(_url).
		Set("Content-Type", "application/json;charset=utf-8").
		Send(payload).
		End()
	if len(es) != 0 || r.StatusCode/100 != 2 {
		slog.Error("发起请求错误", "err", es, "url", _url)
		return ""
	}
	return b
}
