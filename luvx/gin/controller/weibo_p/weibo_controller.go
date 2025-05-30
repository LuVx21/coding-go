package weibo_p

import (
	"net/http"
	"strings"
	"time"

	"luvx/gin/common/responsex"
	"luvx/gin/service/weibo_p"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

func PullByUser(c *gin.Context) {
	uid := c.Query("uid")
	if "0" == uid {
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
	groupIdStr, word, dayStr := c.Query("groupId"), c.Query("word"), c.Query("day")

	groupId := cast_x.ToInt64(groupIdStr)
	var day time.Time
	if len(dayStr) > 0 {
		day, _ = time.Parse(time.DateOnly, dayStr)
	}

	args := map[string]any{
		"groupId": groupId,
		"word":    word,
		"day":     day,

		"size":         c.Query("size"),
		"deleteBefore": c.Query("deleteBefore"),
		"pullBefore":   c.Query("pullBefore"),
	}

	uids := slices_x.Transfer(func(i string) int64 { return cast_x.ToInt64(i) }, strings.Split(uidStr, ",")...)

	rss := weibo_p.Rss(args, groupId, word, day, uids...)
	c.Header("Content-Type", "application/xml;charset=UTF-8")
	c.String(http.StatusOK, rss)
}
func DeleteById(c *gin.Context) {
	id := c.Param("id")
	cnt := weibo_p.DeleteById(cast_x.ToInt64(id))
	responsex.R(c, cnt)
}
