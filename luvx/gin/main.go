package gin

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/config"
    "luvx/gin/middleware"
    "luvx/gin/router"
    "luvx/gin/runner"
)

func WebStart() {
    logs.Log.Infoln("(*^▽^*) 启动... (〃'▽'〃)")
    runner.Start()

    r := gin.Default()

    router.Register(r)
    middleware.RegisterGlobalMiddlewares(r)

    port := config.AppConfig.Server.Port
    if err := r.Run(port); err != nil {
        fmt.Printf("Start server error,err=%v", err)
    }
}
