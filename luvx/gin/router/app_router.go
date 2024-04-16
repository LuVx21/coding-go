package router

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/controller"
    "net/http"
)

type Login struct {
    User     string `form:"user" json:"user" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}

func RegisterApp(r *gin.Engine) {
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "code": "0",
            "msg":  "成功",
            "data": "ok!",
        })
    })

    app := r.Group("/app")
    {
        app.GET("/:path", func(c *gin.Context) {
            path := c.Param("path")
            a := c.Query("a")
            b := c.DefaultQuery("b", "aaa")
            c.JSON(http.StatusOK, gin.H{
                "path": path,
                "a":    a,
                "b":    b,
            })
        })
    }
    {
        app.GET("/healthyCheck", controller.HealthyCheck)
    }

    // 绑定JSON ({"user": "foo", "password": "bar"})
    // 绑定QueryString (/login?user=foo&password=bar)
    r.GET("/login", func(c *gin.Context) {
        var login Login
        if err := c.ShouldBind(&login); err == nil {
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })
}
