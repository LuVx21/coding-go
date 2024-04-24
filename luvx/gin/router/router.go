package router

import (
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/common/consts"
    "luvx/gin/common/responsex"
    "luvx/gin/controller"
    "luvx/gin/controller/weibo_p"
)

// AddTraceId TODO 不太正确
func AddTraceId(c *gin.Context) {
    logs.Log.AddHook(logs.NewTraceIdHook(consts.UUID()))
}

func Register(r *gin.Engine) {
    r.NoMethod(responsex.NoMethod)
    r.NoRoute(responsex.NoRoute)

    routers := []func(*gin.Engine){
        Register0,
        RegisterApp,
        RegisterUser,
        RegisterBili,
        RegisterWeibo,
    }
    for _, router := range routers {
        router(r)
    }
}

func Register0(r *gin.Engine) {
    r.GET("/redirect", controller.Redirect)

    r.GET("/", func(c *gin.Context) {
        logs.Log.Infoln("path:", c.Request.URL.Path)
        responsex.R(c, "ok!")
    })

    app := r.Group("/app")
    app.GET("/healthyCheck", controller.HealthyCheck)
    app.POST("/syncCookie2Turso", controller.SyncCookie2Turso)
}

func RegisterUser(r *gin.Engine) {
    user := r.Group("/user", AddTraceId)
    user.GET("/:username", controller.GetUserByUsername)
}

func RegisterBili(r *gin.Engine) {
    bili := r.Group("/bili", AddTraceId)
    bili.GET("/pull/season", controller.PullSeason)
    bili.GET("/pull/up/video", controller.PullUpVideo)
}

func RegisterWeibo(r *gin.Engine) {
    weibo := r.Group("/weibo")
    weibo.GET("/pull/group", weibo_p.PullByGroup)
    weibo.GET("/pull/user", weibo_p.PullByUser)
    weibo.GET("/rss/:uid", weibo_p.Rss)
    weibo.GET("/rss/delete/:id", weibo_p.DeleteById)
}
