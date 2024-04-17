package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/cast"
    "luvx/gin/service"
    "net/http"
)

func PullSeason(c *gin.Context) {
    seasonId := c.Query("seasonId")
    toInt64 := cast.ToInt64(seasonId)
    if toInt64 == 0 {
        service.PullAll()
    } else {
        service.PullSeasonList(toInt64)
    }
    c.JSON(http.StatusOK, "ok")
}
