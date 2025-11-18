package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/consts_x"
	"luvx/gin/common/consts"
)

func traceId(ctx *gin.Context) {
	ctx.Set(consts_x.TraceId, consts.UUID())
	ctx.Next()
}
