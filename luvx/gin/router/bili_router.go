package router

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller"
)

func RegisterBili(r *gin.Engine) {
    user := r.Group("/bili")
    {
        user.GET("/pull/season", controller.PullSeason)
    }
    {
        user.GET("/pull/up/video", controller.PullUpVideo)
    }
}
