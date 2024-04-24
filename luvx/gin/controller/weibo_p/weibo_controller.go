package weibo_p

import (
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "luvx/gin/common/responsex"
    "luvx/gin/service/weibo_p"
    "net/http"
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
    weibo_p.PullByGroup()
    responsex.R(c, "ok")
}

func Rss(c *gin.Context) {
    uid := c.Param("uid")
    rss := weibo_p.Rss(cast_x.ToInt64(uid))
    c.Header("Content-Type", "application/xml;charset=UTF-8")
    c.String(http.StatusOK, rss)
}
func DeleteById(c *gin.Context) {
    id := c.Param("id")
    cnt := weibo_p.DeleteById(cast_x.ToInt64(id))
    responsex.R(c, cnt)
}
