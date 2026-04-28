package weibo_p

import (
	"net/http"
	"strings"
	"time"

	"luvx/gin/common/responsex"
	"luvx/gin/service/weibo_p"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

func PullByUser(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "0" {
		weibo_p.PullByUserAll()
	} else {
		weibo_p.PullByUser(cast_x.ToInt64(uid))
	}
	responsex.R(c, "ok")
}

func PullByGroup(c *gin.Context) {
	weibo_p.PullByGroupLock()
	weibo_p.DeleteLock()
	responsex.R(c, "ok")
}

func Rss(c *gin.Context) {
	uidStr := c.Param("uid")

	args := map[string]any{"groupId": 0}
	maps_x.ForEach(c.Request.URL.Query(), func(k string, vs []string) {
		if len(vs) >= 1 {
			args[k] = vs[0]
		}
	})

	groupIdStr, dayStr := c.Query("groupId"), c.Query("day")
	if len(groupIdStr) > 0 {
		args["groupId"] = cast_x.ToInt64(groupIdStr)
	}
	var day time.Time
	if len(dayStr) > 0 {
		day, _ = time.Parse(time.DateOnly, dayStr)
	}
	args["day"] = day

	uids := slices_x.Transfer(func(i string) int64 { return cast_x.ToInt64(i) }, strings.Split(uidStr, ",")...)

	rss := weibo_p.Rss(c, args, day, uids...)
	c.Header("Content-Type", "application/xml;charset=UTF-8")
	c.String(http.StatusOK, rss)
}
