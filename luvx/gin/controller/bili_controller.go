package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/cast"
    "luvx/gin/common/responsex"
    "luvx/gin/service/bili"
    "net/http"
)

func PullSeason(c *gin.Context) {
    seasonId := c.Query("seasonId")
    toInt64 := cast.ToInt64(seasonId)
    if toInt64 == 0 {
        bili.PullAll()
    } else {
        bili.PullSeasonList(toInt64)
    }
    c.JSON(http.StatusOK, "ok")
}

func PullUpVideo(c *gin.Context) {
    mid := c.Query("mid")
    video := bili.PullUpVideo(mid)
    responsex.Result(c, video)
}
