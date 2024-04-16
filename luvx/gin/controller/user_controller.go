package controller

import (
    "github.com/gin-gonic/gin"
    "luvx/gin/service"
    "net/http"
)

func GetUserByUsername(c *gin.Context) {
    username := c.Param("username")
    user, err := service.GetUserByUsername(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}
