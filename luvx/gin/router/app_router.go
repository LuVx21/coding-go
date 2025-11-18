package router

import (
	"github.com/gin-gonic/gin"
	"luvx/gin/common/errorx"
	"luvx/gin/common/responsex"
	"net/http"
)

type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func RegisterApp(r *gin.Engine) {
	// 绑定JSON ({"user": "foo", "password": "bar"})
	// 绑定QueryString (/login?user=foo&password=bar)
	r.GET("/login", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err == nil {
			responsex.R(c, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			responsex.R(c, errorx.NewCodeMsgError(http.StatusBadRequest, "异常"))
		}
	})
}
