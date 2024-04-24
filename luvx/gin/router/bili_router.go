package router

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller"
)

func RegisterBili(r *gin.Engine) {
    bili := r.Group("/bili", AddTraceId)
    bili.GET("/pull/season", controller.PullSeason)
    bili.GET("/pull/up/video", controller.PullUpVideo)
}
