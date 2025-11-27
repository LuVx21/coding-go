package controller

import (
	"luvx/gin/common/responsex"
	"luvx/gin/service/bili"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
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
