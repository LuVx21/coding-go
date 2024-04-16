package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func traceId() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        traceId := uuid.New().String()
        ctx.Set("traceId", traceId)
        ctx.Next()
    }
}
