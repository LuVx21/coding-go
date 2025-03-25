package deepseek

import (
	"fmt"
	"os"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/infra/ai"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

const (
	MODELS_URL = "https://%s/models"
)

var (
	aiKey = common_x.IfThen(len(ai.AI_KEY) != 0, ai.AI_KEY, os.Getenv("AI_KEY"))
)

func Models(domain string) []string {
	_, r, _ := gorequest.New().Timeout(time.Minute).
		Get(fmt.Sprintf(MODELS_URL, domain)).
		Set("Content-Type", "application/json").
		Set("Authorization", "Bearer "+aiKey).
		End()

	s := gjson.Get(r, "data.#.id").Array()
	return slices_x.Transfer(func(i gjson.Result) string { return i.String() }, s...)
}
