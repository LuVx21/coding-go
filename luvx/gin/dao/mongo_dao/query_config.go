package mongo_dao

import (
	"context"
	"luvx/gin/db"

	"github.com/icloudza/fxjson"
	"github.com/luvx21/coding-go/coding-common/func_x"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DynamicConfig = func_x.Lazy(func() bson.M {
		var m bson.M
		e := db.GetCollection("config").FindOne(context.TODO(), bson.M{"_id": "app_config"}).Decode(&m)
		if e != nil {
			log.Warnln("lazy加载异常", e)
		}
		return m
	})
	DynamicCache = func_x.Lazy(func() bson.M {
		var m bson.M
		e := db.GetCollection("config").FindOne(context.TODO(), bson.M{"_id": "app_cache"}).Decode(&m)
		if e != nil {
			log.Warnln("lazy加载异常", e)
		}
		return m
	})
	DynamicSwitch = func_x.Lazy(func() fxjson.Node {
		sr := db.GetCollection("config").FindOne(context.TODO(), bson.M{"_id": "app_switch"})
		a, _ := sr.Raw()
		return fxjson.FromString(a.String())
	})
)
