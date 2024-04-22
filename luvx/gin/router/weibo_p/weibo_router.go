package weibo_p

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller/weibo_p"
)

func RegisterWeibo(r *gin.Engine) {
    user := r.Group("/weibo")
    {
        user.GET("/pull/user", weibo_p.PullByUser)
    }
    {
        user.GET("/pull/group", weibo_p.PullByGroup)
    }
}
