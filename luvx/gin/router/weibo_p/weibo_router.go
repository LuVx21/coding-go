package weibo_p

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller/weibo_p"
)

func RegisterWeibo(r *gin.Engine) {
    weibo := r.Group("/weibo")
    weibo.GET("/pull/user", weibo_p.PullByUser)
    weibo.GET("/pull/group", weibo_p.PullByGroup)
}
