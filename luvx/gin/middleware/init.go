package middleware

import (
    "github.com/gin-gonic/gin"
)

func RegisterGlobalMiddlewares(r *gin.Engine) {
    r.Use(auth())
    r.Use(traceId())
    r.Use(requestLog())
    r.Use(recoverLog())
}
