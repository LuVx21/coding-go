package middleware

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/luvx21/coding-go/infra/logs"

	"github.com/gin-gonic/gin"
)

func requestLog(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	if raw != "" {
		path = path + "?" + raw
	}
	Cost := time.Since(start)
	if Cost > time.Minute {
		Cost = Cost.Truncate(time.Second)
	}

	requestMap := map[string]any{
		"Path":      path,
		"Method":    c.Request.Method,
		"ClientIP":  c.ClientIP(),
		"Cost":      Cost.String(),
		"Status":    c.Writer.Status(),
		"Proto":     c.Request.Proto,
		"UserAgent": c.Request.UserAgent(),
		"Msg":       c.Errors.ByType(gin.ErrorTypePrivate).String(),
		"Size":      c.Writer.Size(),
	}

	// logx.WithContext(c).Serve(requestMap)
	j, _ := sonic.Marshal(requestMap)
	logs.Log.Warnln(string(j))
}
