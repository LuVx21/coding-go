package controller

import (
	"luvx/gin/common/responsex"
	"luvx/gin/runner"

	"github.com/gin-gonic/gin"
)

func CallRunner(c *gin.Context) {
	_json := make(map[string]any)
	_ = c.BindJSON(&_json)

	runner.RunnerMap[_json["name"].(string)]()
	responsex.R(c, "ok")
}
