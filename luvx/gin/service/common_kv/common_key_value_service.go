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

func GetByCursor(cursorID int, limit int, bizType common_kv_dao.CommonKVBizType, keys ...string) ([]*model.CommonKeyValue, int, error) {
	return common_kv_dao.GetByCursor(cursorID, limit, int32(bizType), keys...)
}

func Create(kv *model.CommonKeyValue) error {
	return common_kv_dao.Create(kv)
}

func Delete(ids []int) error {
	return common_kv_dao.Delete(ids)
}

func Update(kv *model.CommonKeyValue) error {
	return common_kv_dao.Update(kv)
}
