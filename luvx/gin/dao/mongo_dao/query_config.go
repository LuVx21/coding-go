package mongo_dao

import (
	"context"

	"github.com/icloudza/fxjson"
	"github.com/luvx21/coding-go/coding-common/func_x"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DynamicConfig = func_x.Lazy(func() bson.M {
		var m bson.M
		e := ConfigCol.FindOne(context.TODO(), bson.M{"key": "app_config"}).Decode(&m)
		if e != nil {
			log.Warnln("lazy加载异常", e)
		}
		return m
	})
	DynamicCache = func_x.Lazy(func() bson.M {
		var m bson.M
		e := ConfigCol.FindOne(context.TODO(), bson.M{"key": "app_cache"}).Decode(&m)
		if e != nil {
			log.Warnln("lazy加载异常", e)
		}
		return m
	})
	_dynamicSwitch = func_x.Lazy(func() fxjson.Node {
		sr := ConfigCol.FindOne(context.TODO(), bson.M{"key": "app_switch"})
		a, _ := sr.Raw()
		return fxjson.FromString(a.String())
	})
	DynamicSwitch = func(key string) bool { return _dynamicSwitch.Get().Get(key).BoolOr(false) }
)

var (
	BiliSeason = func() bson.M {
		var result bson.M
		ConfigCol.FindOne(context.TODO(), bson.M{"key": "bili_season"}).Decode(&result)
		return result
	}
)
