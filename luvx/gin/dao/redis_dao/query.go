package redis_dao

import (
	"context"

	"luvx/gin/common/consts"
	"luvx/gin/db"

	"github.com/sirupsen/logrus"
)

const (
	R_K_REMOTE_COOKIE = "remote_cookie"
)

func GetSwitch(key string) bool {
	result, err := db.RedisClient.HGet(context.TODO(), consts.AppSwitchKey, key).Bool()
	if err != nil {
		logrus.Warnln("开关未启用", err, result, "redis key="+key)
	}
	return err == nil && result
}

func SetSwitch(key string, value bool) bool {
	r := db.RedisClient.HSet(context.TODO(), consts.AppSwitchKey, key, value)
	return r.Val() > 0
}
