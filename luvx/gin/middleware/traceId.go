package middleware

import (
	"luvx/gin/common/consts"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/consts_x"
)

func traceId(ctx *gin.Context) {
	ctx.Set(consts_x.TraceId, consts.UUID())
	ctx.Next()
}
