package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"luvx/gin/common/consts"
	"luvx/gin/common/responsex"
	"luvx/gin/dao"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/db"
	"luvx/gin/model"
	"luvx/gin/service/cookie"
	"luvx/gin/service/rpc"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/luvx21/coding-go/coding-common/common_x/t"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	dbs "github.com/luvx21/coding-go/infra/infra_sql"
	"github.com/luvx21/coding-go/luvx_service_sdk/proto_gen/proto_kv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/time/rate"
)

var (
	redisClient = db.InitRedisClient()
)

func HealthyCheck(c *gin.Context) {
	args := 1

	mysql := func() any {
		return common_x.RunWithTimeReturn("mysql", func() model.User {
			return *(dao.GetUserById(args))
		})
	}
	mongo := func() any {
		return common_x.RunWithTimeReturn("mongodb", func() t.Pair[bson.M, error] {
			var result bson.M
			a := mongo_dao.UserCol.FindOne(context.TODO(), bson.M{"_id": args}).Decode(&result)
			return t.NewPair(result, a)
		}).K
	}
	redis := func() any {
		return common_x.RunWithTimeReturn("redis", func() t.Pair[string, error] {
			return t.NewPair(redisClient.Get(context.Background(), "foo").Result())
		}).K
	}
	sqlite := func() any {
		return common_x.RunWithTimeReturn("sqlite", func() t.Pair[[]map[string]any, error] {
			return t.NewPair(dbs.RowsMap(context.TODO(), db.SqliteClient, "select * from user where id = ?", args))
		}).K
	}
	cookie := func() any {
		return common_x.RunWithTimeReturn("cookie", func() a.Table[string] {
			return cookie.GetCookieFromDb(".weibo.com", "weibo.com")
		})
	}
	turso := func() any {
		return common_x.RunWithTimeReturn("turso", func() t.Pair[[]map[string]any, error] {
			return t.NewPair(dbs.RowsMap(context.TODO(), db.Turso, "select * from user where id = ?", args))
		}).K
	}
	fs := []func() any{mysql, mongo, redis, sqlite, cookie, turso}
	rs := make([]any, len(fs))

	wg := sync.WaitGroup{}
	for i, f := range fs {
		wg.Go(func() { rs[i] = f() })
	}
	wg.Wait()

	responsex.R(c, gin.H{
		"mysql":  rs[0],
		"mongo":  rs[1],
		"redis":  rs[2],
		"sqlite": rs[3],
		"cookie": rs[4],
		"turso":  rs[5],
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
			if ignoreHeaders.Contains(k) {
				continue
			}
			c.Header(k, v[0])
		}
		c.String(response.StatusCode, body)
	}
}

func SyncCookie2Yun(c *gin.Context) {
	_json := make(map[string]any)
	_ = c.BindJSON(&_json)
	hosts := slices_x.Transfer(func(a any) string { return a.(string) }, _json["hosts"].([]any)...)
	cookie.Sync2Yun(hosts...)
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
