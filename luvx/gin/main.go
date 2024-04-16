package gin

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/config"
    _ "luvx/gin/db"
    "luvx/gin/middleware"
    "luvx/gin/router"
)

func WebStart() {
    logs.Log.Infoln("(*^▽^*) 启动... (〃'▽'〃)")
    r := gin.Default()

    router.Register(r)
    middleware.RegisterGlobalMiddlewares(r)

    port := config.AppConfig.Server.Port
    if err := r.Run(port); err != nil {
        fmt.Printf("Start server error,err=%v", err)
    }
}
