package redis_dao

import (
	"context"
	"luvx/gin/common/consts"
	"luvx/gin/db"

	"github.com/sirupsen/logrus"
)

func GetSwitch(key string) bool {
	result, err := db.RedisClient.HGet(context.TODO(), consts.AppSwitchKey, key).Bool()
	if err != nil {
		logrus.Warn("定时任务未启用", err, result, "redis key="+key)
	}
	return err == nil && result
}
