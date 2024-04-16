package router

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller"
)

func RegisterUser(r *gin.Engine) {
    user := r.Group("/user")
    {
        user.GET("/:username", controller.GetUserByUsername)
    }
}
