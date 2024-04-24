package middleware

import (
    "github.com/gin-gonic/gin"
)

func RegisterGlobalMiddlewares(r *gin.Engine) {
    middlewares := []gin.HandlerFunc{
        auth,
        traceId,
        requestLog,
        recoverLog,
    }
    for _, middleware := range middlewares {
        r.Use(middleware)
    }
}
