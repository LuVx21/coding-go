package router

import (
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/common/consts"
    "luvx/gin/common/responsex"
    "luvx/gin/router/weibo_p"
)

// AddTraceId TODO 不太正确
func AddTraceId(c *gin.Context) {
    logs.Log.AddHook(logs.NewTraceIdHook(consts.UUID()))
}

func Register(r *gin.Engine) {
    r.NoMethod(func(ctx *gin.Context) {
        responsex.NoMethod(ctx)
        return
    })
    r.NoRoute(func(ctx *gin.Context) {
        responsex.NoRoute(ctx)
        return
    })

    RegisterApp(r)
    RegisterUser(r)
    RegisterBili(r)
    weibo_p.RegisterWeibo(r)
}
