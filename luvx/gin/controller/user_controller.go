package controller

import (
	"net/http"

	"luvx/gin/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetUserByUsername(c *gin.Context) {
	log.Infoln("path:", c.Request.URL.Path)
	username := c.Param("username")
	user, err := service.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
