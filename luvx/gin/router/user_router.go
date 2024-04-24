package router

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller"
)

func RegisterUser(r *gin.Engine) {
    user := r.Group("/user", AddTraceId)
    user.GET("/:username", controller.GetUserByUsername)
}
