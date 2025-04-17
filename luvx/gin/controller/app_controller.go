package controller

import (
	"context"
	"luvx_service_sdk/proto_gen/proto_kv"
	"sync"

	"luvx/gin/common/consts"
	"luvx/gin/common/responsex"
	"luvx/gin/db"
	"luvx/gin/model"
	"luvx/gin/service/cookie"
	"luvx/gin/service/rpc"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	dbs "github.com/luvx21/coding-go/infra/infra_sql"
	"go.mongodb.org/mongo-driver/bson"
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
			userTable := db.MongoDatabase.Collection("user")
			filter := bson.D{{Key: "_id", Value: args}}
			var result bson.M
			a := userTable.FindOne(context.TODO(), filter).Decode(&result)
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
	common_x.RunInRoutine(&wg, func() { r0 <- f0() })
	common_x.RunInRoutine(&wg, func() { r1 <- f1() })
	common_x.RunInRoutine(&wg, func() { r2 <- f2() })

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


func Redirect(c *gin.Context) {
	toUrl := c.Query("url")
	// logs.Log.Infoln("重定向到:", toUrl)

	response, body, _ := consts.GoRequest.
		Proxy("http://" + consts.ServiceHost + ":7890").
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
	gr, _ := rpc.KvRpcClient.Get(context.Background(), &proto_kv.Key{Key: c.Query("key")})
	responsex.R(c, string(gr.Value))
}
