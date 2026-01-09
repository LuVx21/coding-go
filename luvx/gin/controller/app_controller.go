package controller

import (
	"context"
	"fmt"
	"luvx_service_sdk/proto_gen/proto_kv"
	"net/http"
	"strings"
	"sync"

	"luvx/gin/common/consts"
	"luvx/gin/common/responsex"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/db"
	"luvx/gin/model"
	"luvx/gin/service/cookie"
	"luvx/gin/service/rpc"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	dbs "github.com/luvx21/coding-go/infra/infra_sql"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/time/rate"
)

var (
	redisClient = db.InitRedisClient()
)

func HealthyCheck(c *gin.Context) {
	args := 1

	f0 := func() model.User {
		mysql := common_x.RunWithTime("mysql", func() model.User {
			var user model.User
			if err := db.MySQLClient.Where("id = ?", args).First(&user).Error; err != nil {
				panic(err)
			}
			return user
		})
		return mysql
	}

	f1 := func() bson.M {
		mongo, _ := common_x.RunWithTime2("mongodb", func() (bson.M, error) {
			var result bson.M
			a := mongo_dao.UserCol.FindOne(context.TODO(), bson.M{"_id": args}).Decode(&result)
			return result, a
		})
		return mongo
	}
	f2 := func() string {
		redis, _ := common_x.RunWithTime2("redis", func() (string, error) {
			return redisClient.Get(context.Background(), "foo").Result()
		})
		return redis
	}

	wg := sync.WaitGroup{}
	r0 := make(chan model.User, 1)
	r1 := make(chan bson.M, 1)
	r2 := make(chan string, 1)
	wg.Go(func() { r0 <- f0() })
	wg.Go(func() { r1 <- f1() })
	wg.Go(func() { r2 <- f2() })

	sqlite, _ := common_x.RunWithTime2("sqlite", func() ([]map[string]any, error) {
		return dbs.RowsMap(context.TODO(), db.SqliteClient, "select * from user where id = ?", args)
	})
	cookie := common_x.RunWithTime("cookie", func() map[string]string {
		return cookie.GetCookieByHost(".weibo.com", "weibo.com")
	})
	turso, _ := common_x.RunWithTime2("turso", func() ([]map[string]any, error) {
		return dbs.RowsMap(context.TODO(), db.Turso, "select * from user where id = ?", args)
	})

	wg.Wait()
	close(r0)
	close(r1)
	close(r2)
	mysql := <-r0
	mongo := <-r1
	redis := <-r2
	responsex.R(c, gin.H{
		"mysql":  mysql,
		"mongo":  mongo,
		"redis":  redis,
		"sqlite": sqlite,
		"cookie": cookie,
		"turso":  turso,
	})
}

var ignoreHeaders = sets.NewSet("x-frame-options", "X-Frame-Options")

// var limiter = rate.NewLimiter(3, 3)
var limiter = rate.NewLimiter(
	rate.Limit(cast_x.ToFloat64(mongo_dao.DynamicConfig.Get()["redirect_limit"])),
	cast_x.ToInt(mongo_dao.DynamicConfig.Get()["redirect_burst"]),
)

func Redirect(c *gin.Context) {
	toUrl := c.Query("url")
	limiter.Wait(context.TODO())
	// logrus.Infoln("重定向到:", toUrl)
	if rpc.KvRpcClient != nil && strings.Contains(toUrl, ".sinaimg.cn") {
		gr, _ := (*rpc.KvRpcClient).Get(context.Background(), &proto_kv.Key{Key: toUrl})
		if gr != nil && len(gr.Value) > 0 {
			c.Data(http.StatusOK, "image/jpeg", gr.Value)
			return
		}
	}
	fmt.Println("再次请求:", toUrl)
	response, body, _ := consts.GoRequest.
		Proxy("http://" + consts.AppProxy).
		Get(toUrl).
		End()
	if response != nil {
		for k, v := range response.Header {
			if ignoreHeaders.Contain(k) {
				continue
			}
			c.Header(k, v[0])
		}
		c.String(response.StatusCode, body)
	}
}

func SyncCookie2Turso(c *gin.Context) {
	_json := make(map[string]any)
	_ = c.BindJSON(&_json)
	hosts := slices_x.Transfer(func(a any) string { return a.(string) }, _json["hosts"].([]any)...)
	cookie.Sync2Turso(hosts...)
	responsex.R(c, hosts)
}

func ClearCache(c *gin.Context) {
	_ = cookie.ClearCache()

	responsex.R(c, "ok")
}

func KvSet(c *gin.Context) {}
func KvGet(c *gin.Context) {
	gr, _ := (*rpc.KvRpcClient).Get(context.Background(), &proto_kv.Key{Key: c.Query("key")})
	responsex.R(c, string(gr.Value))
}
