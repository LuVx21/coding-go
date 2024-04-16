package middleware

import (
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/common/responsex"

    "github.com/gin-gonic/gin"
)

func recoverLog() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                logs.Log.Error(err)
                responsex.ServiceUnavailable(ctx)
                ctx.Abort()
                return
            }
        }()
        ctx.Next()
    }
}
