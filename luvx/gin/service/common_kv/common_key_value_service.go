package common_kv

import (
	"luvx/gin/dao/common_kv_dao"
	"luvx/gin/model"
)

func Get(bizType common_kv_dao.CommonKVBizType, keys ...string) map[string]*model.CommonKeyValue {
	kvs := common_kv_dao.Get(int32(bizType), keys...)
	m := make(map[string]*model.CommonKeyValue)
	for _, kv := range kvs {
		m[kv.CommonKey] = kv
	}
	return m
}
