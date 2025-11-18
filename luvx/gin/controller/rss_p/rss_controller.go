package rss_p

import (
	"github.com/gin-gonic/gin"
	"luvx/gin/service/rss"
	"net/http"
)

func Rss(c *gin.Context) {
	rss := rss.Rss(c.Param("spiderKey"))
	c.Header("Content-Type", "application/xml;charset=UTF-8")
	c.String(http.StatusOK, rss)
}

func PullByKey(c *gin.Context) {
	rss.PullByKey()
	c.String(http.StatusOK, "ok")
}
