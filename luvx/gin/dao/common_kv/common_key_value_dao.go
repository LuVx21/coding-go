package common_kv

import (
    "luvx/gin/db"
    "luvx/gin/model"
)

func Get(bizType int32, keys ...string) []*model.CommonKeyValue {
    client := db.MySQLClient
    var kvs []*model.CommonKeyValue
    _ = client.Debug().Where("biz_type = ? and invalid = 0 and common_key in ?", bizType, keys).Find(&kvs).Error
    return kvs
}
