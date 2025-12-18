package router

import (
	"luvx/gin/common/consts"
	"luvx/gin/common/responsex"
	"luvx/gin/controller"
	"luvx/gin/controller/ai_c"
	"luvx/gin/controller/common_kv_controller"
	"luvx/gin/controller/rss_p"
	"luvx/gin/controller/useful_c"
	"luvx/gin/controller/weibo_p"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/infra/logs"
	log "github.com/sirupsen/logrus"
)

// AddTraceId TODO 不太正确
func AddTraceId(c *gin.Context) {
	log.AddHook(logs.NewTraceIdHook(consts.UUID()))
}

func Register(r *gin.Engine) {
	r.NoMethod(responsex.NoMethod)
	r.NoRoute(responsex.NoRoute)

	routers := []func(*gin.Engine){
		Register0,
		RegisterApp,
		RegisterUser,
		RegisterBili,
		RegisterWeibo,
	}
	for _, router := range routers {
		router(r)
	}
}

func Register0(r *gin.Engine) {
	r.GET("/redirect", controller.Redirect)

	r.GET("/", func(c *gin.Context) {
		log.Infoln("path:", c.Request.URL.Path)
		responsex.R(c, "ok!")
	})

	app := r.Group("/app")
	app.GET("/healthyCheck", controller.HealthyCheck)
	app.POST("/runner", controller.CallRunner)

	cookie := r.Group("/cookie")
	cookie.POST("/syncCookie2Turso", controller.SyncCookie2Turso)

	useful := r.Group("/useful")
	useful.POST("/compare", useful_c.Compare)
	useful.GET("/c_k_v", common_kv_controller.GetCommonKeyValue)
	useful.POST("/c_k_v", common_kv_controller.CreateCommonKeyValue)
	useful.DELETE("/c_k_v", common_kv_controller.DeleteCommonKeyValue)
	useful.PUT("/c_k_v/:id", common_kv_controller.UpdateCommonKeyValue)

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
	_rss.GET("pullbykey", rss_p.PullByKey)

	_ai := r.Group("/ai")
	_ai.POST("/v1/chat/completions", ai_c.HandleChatCompletion)
}
