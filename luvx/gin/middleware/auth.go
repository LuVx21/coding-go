package middleware

import (
    "github.com/gin-gonic/gin"
)

func auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: 验证身份认证信息
        c.Next()
    }
}
