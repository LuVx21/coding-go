package controller

import (
	"luvx/gin/common/responsex"
	"luvx/gin/service/bili"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

func PullSeason(c *gin.Context) {
	seasonId := c.Query("seasonId")
	toInt64 := cast_x.ToInt64(seasonId)
	if toInt64 == 0 {
		bili.PullAllSeason()
	} else {
		bili.PullSeasonList(toInt64)
	}
	c.JSON(http.StatusOK, "ok")
}

func PullUpVideo(c *gin.Context) {
	mid := c.Query("mid")
	toInt64 := cast_x.ToInt64(mid)
	var video []string
	if toInt64 == 0 {
		bili.PullAllUpVideo()
	} else {
		video = bili.PullUpVideo(toInt64)
	}
	responsex.R(c, video)
}

func Rss(c *gin.Context) {
	uname, includeUid, excludeUid, size := c.Query("uname"), c.Query("includeUid"), c.Query("excludeUid"), c.Query("size")
	f := func(ids string) []int64 {
		return slices_x.FilterTransfer(func(s string) bool { return s != "" }, func(s string) int64 { return cast_x.ToInt64(s) }, strings.Split(ids, ",")...)
	}
	includeUids, excludeUids := f(includeUid), f(excludeUid)

	rss := bili.Rss(uname, includeUids, excludeUids, cast_x.ToInt64(size))
	c.Header("Content-Type", "application/xml;charset=UTF-8")
	c.String(http.StatusOK, rss)
}
