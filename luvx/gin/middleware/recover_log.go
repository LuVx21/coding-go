package middleware

import (
	"luvx/gin/common/responsex"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func recoverLog(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			responsex.ServiceUnavailable(ctx)
			ctx.Abort()
			return
		}
	}()
	ctx.Next()
}
