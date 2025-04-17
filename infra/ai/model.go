package ai

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

const (
	MODELS_URL = "https://%s/models"
)

func Models(domain, apiKeyName string) []string {
	_, r, err := gorequest.New().Timeout(time.Minute).
		Get(fmt.Sprintf(MODELS_URL, domain)).
		Set("Content-Type", "application/json").
		Set("Authorization", "Bearer "+os_x.Getenv(apiKeyName)).
		End()
	if err != nil {
		slog.Warn("获取在线模型异常", "err", err)
		return nil
	}

	s := gjson.Get(r, "data.#.id").Array()
	return slices_x.Transfer(func(i gjson.Result) string { return i.String() }, s...)
}
