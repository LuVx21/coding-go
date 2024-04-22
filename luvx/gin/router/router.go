package router

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/common/responsex"
    "luvx/gin/router/weibo_p"
)

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
