package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/service"
    "net/http"
)

func GetUserByUsername(c *gin.Context) {
    logs.Log.Infoln("path:", c.Request.URL.Path)
    username := c.Param("username")
    user, err := service.GetUserByUsername(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}
