package weibo_p

import (
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "luvx/gin/common/responsex"
    "luvx/gin/service/weibo_p"
)

func PullByUser(c *gin.Context) {
    uid := c.Query("uid")
    if "0" == uid {
        weibo_p.PullByUserAll()
    } else {
        weibo_p.PullByUser(cast_x.ToInt64(uid))
    }
    responsex.Result(c, "ok")
}

func PullByGroup(c *gin.Context) {
    weibo_p.PullByGroup()
    responsex.Result(c, "ok")
}
