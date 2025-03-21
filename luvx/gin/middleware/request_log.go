package middleware

import (
	"fmt"
	"time"

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
	TimeStamp := time.Now()
	Cost := TimeStamp.Sub(start)
	if Cost > time.Minute {
		Cost = Cost.Truncate(time.Second)
	}

	requestMap := map[string]interface{}{
		"Path":      path,
		"Method":    c.Request.Method,
		"ClientIP":  c.ClientIP(),
		"Cost":      fmt.Sprintf("%s", Cost),
		"Status":    c.Writer.Status(),
		"Proto":     c.Request.Proto,
		"UserAgent": c.Request.UserAgent(),
		"Msg":       c.Errors.ByType(gin.ErrorTypePrivate).String(),
		"Size":      c.Writer.Size(),
	}

	// logx.WithContext(c).Serve(requestMap)
	logs.Log.Warnln(requestMap)
}
