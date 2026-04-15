package router

import (
	"luvx/gin/common/consts"
	"luvx/gin/common/errorx"
	"luvx/gin/common/responsex"
	"luvx/gin/controller"
	"luvx/gin/controller/ai_c"
	"luvx/gin/controller/ckv_c"
	"luvx/gin/controller/rss_p"
	"luvx/gin/controller/useful_c"
	"luvx/gin/controller/weibo_p"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/infra/logs"
	log "github.com/sirupsen/logrus"
)

type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// AddTraceId TODO 不太正确
func AddTraceId(c *gin.Context) {
	log.AddHook(logs.NewTraceIdHook(consts.UUID()))
}

func Register(r *gin.Engine) {
	r.NoMethod(responsex.NoMethod)
	r.NoRoute(responsex.NoRoute)

	routers := []func(*gin.Engine){
		Register0,
		RegisterUser,
		RegisterBili,
		RegisterWeibo,
	}
	for _, router := range routers {
		router(r)
	}
}

func Register0(r *gin.Engine) {
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

	r.GET("/redirect", controller.Redirect)

	r.GET("/", func(c *gin.Context) {
		log.Infoln("path:", c.Request.URL.Path)
		responsex.R(c, "ok!")
	})

	app := r.Group("/app")
	app.GET("/healthyCheck", controller.HealthyCheck)
	app.POST("/runner", controller.CallRunner)

	cookie := r.Group("/cookie")
	cookie.POST("/syncCookie2Yun", controller.SyncCookie2Yun)

	useful := r.Group("/useful")
	useful.POST("/compare", useful_c.Compare)

	useful.GET("/c_k_v", ckv_c.GetCommonKeyValue)
	useful.POST("/c_k_v", ckv_c.CreateCommonKeyValue)
	useful.DELETE("/c_k_v", ckv_c.DeleteCommonKeyValue)
	useful.PUT("/c_k_v/:id", ckv_c.UpdateCommonKeyValue)

	cache := r.Group("/cache")
	cache.GET("clear", controller.ClearCache)

	kv := r.Group("/kv")
	kv.GET("get", controller.KvGet)
}

func RegisterUser(r *gin.Engine) {
	user := r.Group("/user", AddTraceId)
	user.GET("/:username", controller.GetUserByUsername)
}

func RegisterBili(r *gin.Engine) {
	bili := r.Group("/bili", AddTraceId)
	bili.GET("/pull/season", controller.PullSeason)
	bili.GET("/pull/up/video", controller.PullUpVideo)
	bili.GET("/rss", controller.Rss)
}

func RegisterWeibo(r *gin.Engine) {
	weibo := r.Group("/weibo")
	weibo.GET("/pull/group", weibo_p.PullByGroup)
	weibo.GET("/pull/user", weibo_p.PullByUser)
	weibo.GET("/rss/:uid", weibo_p.Rss)

	_rss := r.Group("/rss")
	_rss.GET("/feed/:spiderKey", rss_p.Rss)
	_rss.GET("/delete/:source/:id", rss_p.DeleteById)
	_rss.GET("/info/pulled", rss_p.RssPull)

	_ai := r.Group("/ai")
	_ai.POST("/v1/chat/completions", ai_c.HandleChatCompletion)
}
